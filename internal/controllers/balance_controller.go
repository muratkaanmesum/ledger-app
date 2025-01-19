package controllers

import (
	"ptm/internal/repositories"
)

type BalanceController struct {
	Repo repositories.BalanceRepository
}

func (b *BalanceController) getBalance(userId uint) {
}
