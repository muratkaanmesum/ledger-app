package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"ptm/internal/db/transaction"
	"ptm/internal/di"
	"ptm/internal/models"
	"ptm/internal/services"
	"ptm/internal/utils/customError"
	"ptm/internal/utils/response"
	"ptm/pkg/jwt"
	"ptm/pkg/worker"
)

type TransactionController interface {
	HandleCredit(c echo.Context) error
	HandleDebit(c echo.Context) error
}

type transactionController struct {
	balanceService     services.BalanceService
	transactionService services.TransactionService
	pool               *worker.Pool
}

func NewTransactionController() TransactionController {
	pool := worker.GetPool()

	return &transactionController{
		balanceService:     di.Resolve[services.BalanceService](),
		transactionService: di.Resolve[services.TransactionService](),
		pool:               pool,
	}
}

type creditRequest struct {
	Amount float64 `json:"amount" validate:"required"`
}

func (tc *transactionController) HandleCredit(c echo.Context) error {
	var (
		req  creditRequest
		user = c.Get("user").(*jwt.CustomClaims)
	)
	if err := c.Bind(&req); err != nil {
		return customError.New(customError.InternalServerError, err)
	}

	if err := c.Validate(req); err != nil {
		fmt.Println("WOW ITS HERE ", err)
		return customError.New(customError.BadRequest, err)
	}

	tc.pool.AddTask(func() {
		fmt.Println("RUNNING GOROUTİNE ")
		db, err := transaction.StartTransaction()
		if err != nil {
			//return customError.New(customError.InternalServerError, err)
		}

		if err := tc.balanceService.IncrementUserBalance(user.Id, req.Amount); err != nil {
			if rollbackErr := transaction.RollbackTransaction(db); rollbackErr != nil {
				//	return customError.New(customError.InternalServerError, rollbackErr)
			}
			// return customError.New(customError.InternalServerError, err)
		}

		createdTransaction, err := tc.transactionService.CreateTransaction(
			user.Id, user.Id, req.Amount, models.Credit,
		)
		if err != nil {
			if rollbackErr := transaction.RollbackTransaction(db); rollbackErr != nil {
				//return customError.New(customError.InternalServerError, rollbackErr)
			}
			//return customError.New(customError.InternalServerError, err)
		}

		if err := transaction.CommitTransaction(db); err != nil {
			// return customError.New(customError.InternalServerError, err)
		}
		fmt.Println(createdTransaction)
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

	db, err := transaction.StartTransaction()
	if err != nil {
		return response.InternalServerError(c, "Error starting transaction", err)
	}

	if err := tc.balanceService.DecrementUserBalance(user.Id, req.Amount); err != nil {
		if rollbackErr := transaction.RollbackTransaction(db); rollbackErr != nil {
			return response.InternalServerError(c, "Error rolling back transaction", rollbackErr)
		}
		return response.BadRequest(c, "Error updating user balance", err)
	}

	createdTransaction, err := tc.transactionService.CreateTransaction(
		user.Id, user.Id, req.Amount, models.Debit,
	)
	if err != nil {
		if rollbackErr := transaction.RollbackTransaction(db); rollbackErr != nil {
			return response.InternalServerError(c, "Error rolling back transaction", rollbackErr)
		}
		return response.InternalServerError(c, "Error creating transaction", err)
	}

	if err := transaction.CommitTransaction(db); err != nil {
		return response.InternalServerError(c, "Error committing transaction", err)
	}

	return response.Ok(c, "Transaction successful", createdTransaction)
}
