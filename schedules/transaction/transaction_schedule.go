package transactionSchedule

import (
	"log"
	"ptm/internal/db/transaction"
	"ptm/internal/di"
	"ptm/internal/models"
	"ptm/internal/repositories"
	"ptm/internal/scheduler"
	"ptm/internal/services"
)

func handleScheduledTransactions() {
	scheduleRepository := di.Resolve[repositories.ScheduleRepository]()

	data, err := scheduleRepository.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, schedule := range data {
		db, err := transaction.StartTransaction()
		if err != nil {
			log.Println("failed to start db transaction", err)
			continue
		}

		createdTransaction, err := di.Resolve[services.TransactionService]().CreateTransaction(
			schedule.ID, schedule.TargetUserID, schedule.Amount, "transfer",
		)
		if err != nil {
			log.Println("failed to create transaction", err)
			_ = transaction.RollbackTransaction(db)
			continue
		}

		if err := di.Resolve[services.BalanceService]().DecrementUserBalance(schedule.ID, schedule.Amount); err != nil {
			log.Println("failed to decrement user balance", err)
			_ = transaction.RollbackTransaction(db)
			continue
		}

		if err := di.Resolve[services.BalanceService]().IncrementUserBalance(schedule.ID, schedule.Amount); err != nil {
			log.Println("failed to increment recipient balance", err)
			_ = transaction.RollbackTransaction(db)
			continue
		}

		if err := di.Resolve[services.TransactionService]().UpdateTransactionState(createdTransaction.ID, models.TransactionStatusCompleted); err != nil {
			log.Println("failed to update transaction state", err)
			_ = transaction.RollbackTransaction(db)
			continue
		}

		if err := transaction.CommitTransaction(db); err != nil {
			log.Println("failed to commit transaction", err)
		}
	}
}
func init() {
	scheduler.AddSchedule(scheduler.CronJob{
		JobFunc: handleScheduledTransactions,
		Spec:    "ScheduledTransactions",
		Time:    "0 5 * * *",
	})
}
