package controllers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ptm/internal/db/transaction"
	"ptm/internal/di"
	"ptm/internal/models"
	"ptm/internal/services"
	"ptm/internal/utils/response"
	"ptm/pkg/jwt"
	"ptm/pkg/logger"
	"ptm/pkg/worker"
)

type TransactionController interface {
	HandleCredit(c echo.Context) error
	HandleDebit(c echo.Context) error
	GetTransactions(c echo.Context) error
	HandleTransfer(c echo.Context) error
}

type transactionController struct {
	balanceService     services.BalanceService
	transactionService services.TransactionService
	userService        services.UserService
	pool               *worker.Pool
}

func NewTransactionController() TransactionController {
	pool := worker.GetPool()

	return &transactionController{
		balanceService:     di.Resolve[services.BalanceService](),
		transactionService: di.Resolve[services.TransactionService](),
		userService:        di.Resolve[services.UserService](),
		pool:               pool,
	}
}

type creditRequest struct {
	Amount float64 `json:"amount" validate:"required"`
}

type TransferRequest struct {
	Amount float64 `json:"amount" validate:"required"`
	ToId   uint    `json:"to_id" validate:"required"`
}

type PaginationRequest struct {
	Page  uint `json:"page"`
	Count uint `json:"count"`
}

func (tc *transactionController) HandleCredit(c echo.Context) error {
	var (
		req  creditRequest
		user = jwt.GetUser(c)
	)

	if err := c.Bind(&req); err != nil {
		return response.InternalServerError(c, "Error binding request", err)
	}

	if err := c.Validate(req); err != nil {
		return response.UnprocessableEntity(c, "Validation error", err)
	}

	tc.pool.AddTask(func() {
		createdTransaction, err := tc.transactionService.CreateTransaction(
			user.Id, user.Id, req.Amount, models.Credit,
		)
		if err != nil {
		}

		db, err := transaction.StartTransaction()
		if err != nil {
			//return customError.New(customError.InternalServerError, err)
		}

		defer tc.HandleTransactionPanic(db, createdTransaction)
		if err := tc.balanceService.IncrementUserBalance(user.Id, req.Amount); err != nil {

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
	})

	return response.Accepted(c, "Transaction Accepted")
}

func (tc *transactionController) HandleDebit(c echo.Context) error {
	var (
		req  creditRequest
		user = jwt.GetUser(c)
	)
	if err := c.Bind(&req); err != nil {
		return response.InternalServerError(c, "Error binding request", err)
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, "Validation error", err)
	}

	tc.pool.AddTask(func() {
		db, err := transaction.StartTransaction()
		if err != nil {
			//return response.InternalServerError(c, "Error starting transaction", err)
		}
		createdTransaction, err := tc.transactionService.CreateTransaction(
			user.Id, user.Id, req.Amount, models.Debit,
		)

		defer tc.HandleTransactionPanic(db, createdTransaction)

		if err := tc.balanceService.DecrementUserBalance(user.Id, req.Amount); err != nil {
			logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", err))
		}

		if err != nil {
			logger.Logger.Error("Error creating transaction", zap.Error(err))
		}

		if err := tc.transactionService.UpdateTransactionState(user.Id, models.TransactionStatusCompleted); err != nil {

		}

		if err := transaction.CommitTransaction(db); err != nil {
			logger.Logger.Error("Error when committing the transaction")
		}
	})

	return response.Accepted(c, "Transaction Accepted")
}

func (tc *transactionController) GetTransactions(c echo.Context) error {
	var (
		req  PaginationRequest
		user = jwt.GetUser(c)
	)

	if err := c.Bind(&req); err != nil {
		return response.InternalServerError(c, "Error binding request", err)
	}

	if err := c.Validate(req); err != nil {
		return response.UnprocessableEntity(c, "Validation error", err)
	}

	transactions, err := tc.transactionService.ListTransactions(user.Id, req.Page, req.Count)

	if err != nil {
		return err
	}

	return response.Ok(c, "Successfully Fetched", transactions)
}

func (tc *transactionController) HandleTransactionPanic(db *gorm.DB, createdTransaction *models.Transaction) {
	if p := recover(); p != nil {
		if err := transaction.RollbackTransaction(db); err != nil {
			logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", p))
			if err := tc.transactionService.UpdateTransactionState(createdTransaction.ID, models.TransactionStatusFailed); err != nil {
				logger.Logger.Error("failed to update transaction status", zap.Error(err))
			}
		}
	}
}

func (tc *transactionController) HandleTransfer(c echo.Context) error {
	user := jwt.GetUser(c)
	var req TransferRequest

	if err := c.Bind(&req); err != nil {
		return response.InternalServerError(c, "Error binding request", err)
	}

	if err := c.Validate(req); err != nil {
		return response.UnprocessableEntity(c, "Validation error", err)
	}

	exists, err := tc.userService.Exists(user.Id)

	if err != nil {
		return err
	}

	if !exists {
		return response.BadRequest(c, "User does not exist")
	}

	tc.pool.AddTask(func() {
		db, err := transaction.StartTransaction()
		if err != nil {
			//return response.InternalServerError(c, "Error starting transaction", err)
		}

		createdTransaction, err := tc.transactionService.CreateTransaction(
			user.Id, req.ToId, req.Amount, models.Transfer,
		)

		defer tc.HandleTransactionPanic(db, createdTransaction)

		if err := tc.balanceService.DecrementUserBalance(user.Id, req.Amount); err != nil {
			logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", err))
		}

		if err := tc.balanceService.IncrementUserBalance(req.ToId, req.Amount); err != nil {
			logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", err))
		}

		if err != nil {
			logger.Logger.Error("Error creating transaction", zap.Error(err))
		}

		if err := tc.transactionService.UpdateTransactionState(user.Id, models.TransactionStatusCompleted); err != nil {

		}

		if err := transaction.CommitTransaction(db); err != nil {
			logger.Logger.Error("Error when committing the transaction")
		}
	})

	return response.Accepted(c, "Transaction Accepted")
}
