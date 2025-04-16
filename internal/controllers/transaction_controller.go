package controllers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ptm/internal/db/transaction"
	"ptm/internal/di"
	"ptm/internal/dtos"
	"ptm/internal/models"
	"ptm/internal/services"
	"ptm/pkg/jwt"
	"ptm/pkg/logger"
	"ptm/pkg/utils/customError"
	"ptm/pkg/utils/response"
	"ptm/pkg/worker"
	"strconv"
	"time"
)

type TransactionController interface {
	HandleCredit(c echo.Context, req CreditRequest) error
	HandleDebit(c echo.Context, req CreditRequest) error
	HandleTransfer(c echo.Context, req TransferRequest) error
	ScheduleTransaction(c echo.Context, req ScheduleRequest) error
	GetTransactions(c echo.Context) error
	GetById(c echo.Context) error
}

type transactionController struct {
	balanceService     services.BalanceService
	transactionService services.TransactionService
	userService        services.UserService
	pool               *worker.Pool
}

func NewTransactionController() TransactionController {
	pool := worker.InitWorkerPool("Transaction", 10)

	return &transactionController{
		balanceService:     di.Resolve[services.BalanceService](),
		transactionService: di.Resolve[services.TransactionService](),
		userService:        di.Resolve[services.UserService](),
		pool:               pool,
	}
}

type CreditRequest struct {
	Amount   float64 `json:"amount" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
}

type TransferRequest struct {
	Amount   float64 `json:"amount" validate:"required"`
	ToId     uint    `json:"to_id" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
}

type ScheduleRequest struct {
	Amount   float64 `json:"amount" validate:"required"`
	ToId     uint    `json:"to_id" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
	time     time.Time
}

func (tc *transactionController) HandleCredit(c echo.Context, req CreditRequest) error {
	var (
		user = jwt.GetUser(c)
	)

	createdTransaction, err := tc.transactionService.CreateTransaction(
		user.Id, user.Id, req.Amount, models.Credit,
	)
	if err != nil {
	}

	db, err := transaction.StartTransaction()
	if err != nil {
		logger.Logger.Error("failed to increment user balance", zap.Error(err))
	}

	defer tc.handleTransactionPanic(db, createdTransaction)
	if err := tc.balanceService.IncrementUserBalance(user.Id, req.Amount, req.Currency); err != nil {
		logger.Logger.Error("failed to increment user balance", zap.Error(err))
	}

	if err != nil {
		logger.Logger.Error("failed to update transaction status", zap.Error(err))
	}

	if err := tc.transactionService.UpdateTransactionState(createdTransaction.ID, models.TransactionStatusCompleted); err != nil {
		logger.Logger.Error("failed to update transaction status", zap.Error(err))
	}

	if err := transaction.CommitTransaction(db); err != nil {
		logger.Logger.Error("failed to commit transaction", zap.Error(err))
	}

	return response.Accepted(c, "Transaction Accepted")
}

func (tc *transactionController) HandleDebit(c echo.Context, req CreditRequest) error {
	var (
		user = jwt.GetUser(c)
	)

	db, err := transaction.StartTransaction()
	if err != nil {
		logger.Logger.Error("failed to increment user balance", zap.Error(err))
	}
	createdTransaction, err := tc.transactionService.CreateTransaction(
		user.Id, user.Id, req.Amount, models.Debit,
	)

	defer tc.handleTransactionPanic(db, createdTransaction)

	if err := tc.balanceService.DecrementUserBalance(user.Id, req.Amount, req.Currency); err != nil {
		logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", err))
	}

	if err != nil {
		logger.Logger.Error("Error creating transaction", zap.Error(err))
	}

	if err := tc.transactionService.UpdateTransactionState(user.Id, models.TransactionStatusCompleted); err != nil {
		logger.Logger.Error("Error creating transaction", zap.Error(err))
	}

	if err := transaction.CommitTransaction(db); err != nil {
		logger.Logger.Error("Error when committing the transaction")
	}

	return response.Accepted(c, "Transaction Accepted")
}

func (tc *transactionController) GetTransactions(c echo.Context) error {
	var (
		req  dtos.PaginationRequest
		user = jwt.GetUser(c)
	)
	if err := c.Bind(&req); err != nil {
		return response.InternalServerError(c, "Error binding request", err)
	}

	failedParam := c.QueryParam("failed")

	if failedParam == "" {
		failedParam = "true"
	}

	failed, err := strconv.ParseBool(failedParam)

	if err != nil {
		return customError.BadRequest("failed", err)
	}
	transactions, err := tc.transactionService.ListTransactions(user.Id, req.Page, req.Count, failed)

	if err != nil {
		return err
	}

	return response.Ok(c, "Successfully Fetched", transactions)
}

func (tc *transactionController) handleTransactionPanic(db *gorm.DB, createdTransaction *models.Transaction) {
	if p := recover(); p != nil {
		if err := transaction.RollbackTransaction(db); err != nil {
			logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", p))
			if err := tc.transactionService.UpdateTransactionState(createdTransaction.ID, models.TransactionStatusFailed); err != nil {
				logger.Logger.Error("failed to update transaction status", zap.Error(err))
			}
		}
	}
}

func (tc *transactionController) HandleTransfer(c echo.Context, req TransferRequest) error {
	user := jwt.GetUser(c)
	exists, err := tc.userService.Exists(user.Id)

	if err != nil {
		return err
	}

	if !exists {
		return response.BadRequest(c, "User does not exist")
	}

	rules, err := tc.userService.GetUserRules(user.Id)
	if err != nil {
		return response.InternalServerError(c, "Error", err)
	}

	if rules.MaxAmountToTransfer < req.Amount {
		return response.BadRequest(c, "Amount is too high", nil)
	}

	db, err := transaction.StartTransaction()
	if err != nil {
		logger.Logger.Error("error", zap.Error(err))
	}

	createdTransaction, err := tc.transactionService.CreateTransaction(
		user.Id, req.ToId, req.Amount, models.Transfer,
	)

	defer tc.handleTransactionPanic(db, createdTransaction)

	if err := tc.balanceService.DecrementUserBalance(user.Id, req.Amount, req.Currency); err != nil {
		logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", err))
	}

	if err := tc.balanceService.IncrementUserBalance(req.ToId, req.Amount, req.Currency); err != nil {
		logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", err))
	}

	if err != nil {
		logger.Logger.Error("Error creating transaction", zap.Error(err))
	}

	if err := tc.transactionService.UpdateTransactionState(user.Id, models.TransactionStatusCompleted); err != nil {
		logger.Logger.Error("failed to update transaction status", zap.Error(err))
	}

	if err := transaction.CommitTransaction(db); err != nil {
		logger.Logger.Error("Error when committing the transaction")
	}

	return response.Accepted(c, "Transaction Accepted")
}

func (tc *transactionController) GetById(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid id")
	}

	returnedTransaction, err := tc.transactionService.GetTransactionById(uint(id))

	if err != nil {
		return err
	}

	return response.Ok(c, "Successful", returnedTransaction)
}

func (tc *transactionController) ScheduleTransaction(c echo.Context, req ScheduleRequest) error {
	user := jwt.GetUser(c)

	if err := c.Bind(&req); err != nil {
		return response.InternalServerError(c, "Error binding request", err)
	}

	if err := c.Validate(req); err != nil {
		return response.UnprocessableEntity(c, "Validation error", err)
	}

	if err := tc.transactionService.ScheduleTransaction(
		user.Id,
		req.ToId,
		req.Amount,
		"credit",
		req.time,
	); err != nil {
		return err
	}

	return response.Accepted(c, "Queued", nil)
}
