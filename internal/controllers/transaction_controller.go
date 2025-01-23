package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ptm/internal/db/transaction"
	"ptm/internal/di"
	"ptm/internal/models"
	"ptm/internal/services"
	"ptm/internal/utils/customError"
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

func (tc *transactionController) HandleCredit(c echo.Context) error {
	var (
		req  creditRequest
		user = jwt.GetUser(c)
	)

	if err := c.Bind(&req); err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	if err := c.Validate(req); err != nil {
		return customError.New(customError.BadRequest, err)
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

		defer HandleTransactionPanic(db)
		if err := tc.balanceService.IncrementUserBalance(user.Id, req.Amount); err != nil {
			if err := tc.transactionService.UpdateTransactionState(createdTransaction.ID, models.TransactionStatusFailed); err != nil {
				logger.Logger.Error("failed to update transaction status", zap.Error(err))
			}
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
		user = c.Get("user").(*jwt.CustomClaims)
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

		defer HandleTransactionPanic(db)

		if err := tc.balanceService.DecrementUserBalance(user.Id, req.Amount); err != nil {
			logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", err))
		}

		_, err = tc.transactionService.CreateTransaction(
			user.Id, user.Id, req.Amount, models.Debit,
		)
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
	user := jwt.GetUser(c)

	transactions, err := tc.transactionService.ListTransactions(user.Id, 10, 0)

	if err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	return response.Ok(c, "Successfully Fetched", transactions)
}

func HandleTransactionPanic(db *gorm.DB) {
	if p := recover(); p != nil {
		if err := transaction.RollbackTransaction(db); err != nil {
			logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", p))
		}
		logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", p))
	} else if err := transaction.RollbackTransaction(db); err != nil {
		logger.Logger.Error("Error occurred during transaction, rolling back", zap.Error(err))
	}
}

func (tc *transactionController) HandleTransfer(c echo.Context) error {
	user := jwt.GetUser(c)
	var req TransferRequest

	if err := c.Bind(&req); err != nil {
		return customError.New(customError.InternalServerError)
	}

	if err := c.Validate(req); err != nil {
		return customError.New(customError.BadRequest, err)
	}

	exists, err := tc.userService.Exists(user.Id)

	if err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	if !exists {
		return customError.New(customError.BadRequest, errors.New("user not found"))
	}

	tc.pool.AddTask(func() {
		db, err := transaction.StartTransaction()
		if err != nil {
			//return response.InternalServerError(c, "Error starting transaction", err)
		}

		defer HandleTransactionPanic(db)

		if err := tc.balanceService.DecrementUserBalance(user.Id, req.Amount); err != nil {
			logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", err))
		}

		if err := tc.balanceService.IncrementUserBalance(req.ToId, req.Amount); err != nil {
			logger.Logger.Error("Panic occurred during transaction, rolling back", zap.Any("panic", err))
		}

		_, err = tc.transactionService.CreateTransaction(
			user.Id, req.ToId, req.Amount, models.Transfer,
		)
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
