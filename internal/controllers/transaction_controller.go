package controllers

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/db/transaction"
	"ptm/internal/models"
	"ptm/internal/repositories"
	"ptm/internal/services"
	"ptm/internal/utils/jwt"
	"ptm/internal/utils/response"
)

type TransactionController struct {
}
type creditRequest struct {
	UserId uint    `json:"sender_id" validate:"required"`
	Amount float64 `json:"amount" validate:"required"`
}

func HandleCredit(c echo.Context) error {
	var (
		req  creditRequest
		user = c.Get("user").(*jwt.CustomClaims)
	)
	if err := c.Bind(&req); err != nil {
		return response.InternalServerError(c, "Error", err)
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, "Error validating body", err)
	}
	db, err := transaction.StartTransaction()

	if err != nil {
		return response.InternalServerError(c, "Error starting transaction", err)
	}

	balanceService := services.NewBalanceService(repositories.NewBalanceRepository())

	if err := balanceService.IncrementUserBalance(req.UserId, req.Amount); err != nil {
		if err := transaction.RollbackTransaction(db); err != nil {
			return response.InternalServerError(c, "Error rolling back transaction", err)
		}

		return response.BadRequest(c, "Error updating user balance", err)
	}

	transactionService := services.NewTransactionService()

	createdTransaction, err := transactionService.CreateTransaction(user.Id, user.Id, req.Amount, models.Credit)
	if err != nil {
		if err := transaction.RollbackTransaction(db); err != nil {
			return response.InternalServerError(c, "Error rolling back transaction", err)
		}

		return response.InternalServerError(c, "Error Creating transaction", err)
	}

	if err := transaction.CommitTransaction(db); err != nil {
		return response.InternalServerError(c, "Error committing transaction", err)
	}

	return response.Ok(c, "Transaction Successful", createdTransaction)
}

func HandleDebit(c echo.Context) error {
	var (
		req  creditRequest
		user = c.Get("user").(*jwt.CustomClaims)
	)
	if err := c.Bind(&req); err != nil {
		return response.InternalServerError(c, "Error", err)
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, "Error validating body", err)
	}
	db, err := transaction.StartTransaction()

	if err != nil {
		return response.InternalServerError(c, "Error starting transaction", err)
	}

	balanceService := services.NewBalanceService(repositories.NewBalanceRepository())

	if err := balanceService.DecrementUserBalance(req.UserId, req.Amount); err != nil {
		if err := transaction.RollbackTransaction(db); err != nil {
			return response.InternalServerError(c, "Error rolling back transaction", err)
		}

		return response.BadRequest(c, "Error updating user balance", err)
	}

	transactionService := services.NewTransactionService()

	createdTransaction, err := transactionService.CreateTransaction(user.Id, user.Id, req.Amount, models.Credit)
	if err != nil {
		if err := transaction.RollbackTransaction(db); err != nil {
			return response.InternalServerError(c, "Error rolling back transaction", err)
		}

		return response.InternalServerError(c, "Error Creating transaction", err)
	}

	if err := transaction.CommitTransaction(db); err != nil {
		return response.InternalServerError(c, "Error committing transaction", err)
	}

	return response.Ok(c, "Transaction Successful", createdTransaction)
}
