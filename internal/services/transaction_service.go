package services

import (
	"fmt"
	"github.com/google/uuid"
	"ptm/internal/event"
	"ptm/internal/models"
	"ptm/internal/repositories"
	"time"
)

type TransactionService interface {
	CreateTransaction(fromId, toId uint, amount float64, transactionType models.TransactionType) (*models.Transaction, error)
	ListTransactions(userID uint, page, count uint, failed bool) ([]models.Transaction, error)
	UpdateTransactionState(transactionId uint, state models.TransactionType) error
	GetTransactionById(transactionId uint) (*models.Transaction, error)
	ScheduleTransaction(
		fromId,
		toId uint,
		amount float64,
		transactionType models.TransactionType,
		execTime time.Time) error
}

type transactionService struct {
	repository         repositories.TransactionRepository
	scheduleRepository repositories.ScheduleRepository
}

func NewTransactionService(repository repositories.TransactionRepository) TransactionService {
	return &transactionService{
		repository: repository,
	}
}

func (t *transactionService) CreateTransaction(fromId, toId uint, amount float64, transactionType models.TransactionType) (*models.Transaction, error) {
	toTransaction, err := models.NewTransaction(fromId, toId, amount, transactionType, models.TransactionStatusPending)
	if err != nil {
		return nil, err
	}

	if err := t.repository.CreateTransaction(toTransaction); err != nil {
		return nil, err
	}

	eventStruct := models.Event{
		ID:        uuid.NewString(),
		EntityID:  fmt.Sprintf("%d", toTransaction.ID),
		Type:      "TransactionCreated",
		Payload:   fmt.Sprintf(`{"from_id":%d,"to_id":%d,"amount":%.2f,"type":"%s"}`, fromId, toId, amount, transactionType),
		Timestamp: time.Now(),
	}
	err = event.AppendEvent("transaction_events", eventStruct)
	if err != nil {
		return nil, err
	}

	return toTransaction, nil
}

func (t *transactionService) UpdateTransactionState(transactionId uint, state models.TransactionType) error {
	transaction, err := t.repository.GetTransactionByID(transactionId)

	if err != nil {
		return err
	}

	transaction.Status = string(state)
	if err := t.repository.UpdateTransaction(transaction); err != nil {
		return err
	}

	return nil
}

func (t *transactionService) ListTransactions(userID uint, page, count uint, failed bool) ([]models.Transaction, error) {
	var transactions []models.Transaction

	transactions, err := t.repository.GetAllTransactions(userID, page, count, failed)

	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (t *transactionService) GetTransactionById(transactionId uint) (*models.Transaction, error) {
	transaction, err := t.repository.GetTransactionByID(transactionId)

	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t *transactionService) ScheduleTransaction(
	fromId,
	toId uint,
	amount float64,
	transactionType models.TransactionType,
	execTime time.Time,
) error {
	return t.scheduleRepository.Create(&models.Schedule{
		Amount:          amount,
		UserID:          fromId,
		TargetUserID:    toId,
		ExecuteAt:       execTime,
		TransactionType: transactionType,
	})
}
