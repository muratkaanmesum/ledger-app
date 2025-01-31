package di

import (
	"fmt"
	"ptm/internal/repositories"
	"ptm/internal/services"
	"reflect"
)

type Container struct {
	singleton map[reflect.Type]reflect.Value
}

var diContainer *Container

func NewContainer() *Container {
	return &Container{
		singleton: make(map[reflect.Type]reflect.Value),
	}
}

func (c *Container) RegisterSingleton(serviceType any, service any) error {
	t := reflect.TypeOf(serviceType).Elem()
	c.singleton[t] = reflect.ValueOf(service)
	return nil
}

func (c *Container) resolve(serviceType any) any {
	t := reflect.TypeOf(serviceType).Elem()

	if resolved, exists := c.singleton[t]; exists {
		return resolved.Interface()
	}

	return nil
}

func Resolve[T any]() T {
	var zero T
	resolved, ok := diContainer.resolve((*T)(nil)).(T)
	if !ok || reflect.ValueOf(resolved).IsNil() {
		panic(fmt.Sprintf("Failed to resolve dependency: %T", zero))
	}
	return resolved
}

func registerRepositories(container *Container) {
	userRepository := repositories.NewUserRepository()
	balanceRepository := repositories.NewBalanceRepository()
	transactionRepository := repositories.NewTransactionRepository()
	balanceHistoryRepository := repositories.NewBalanceHistoryRepository()
	auditRepository := repositories.NewAuditLogRepository()

	if err := container.RegisterSingleton((*repositories.UserRepository)(nil), userRepository); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*repositories.TransactionRepository)(nil), transactionRepository); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*repositories.BalanceRepository)(nil), balanceRepository); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*repositories.BalanceHistoryRepository)(nil), balanceHistoryRepository); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*repositories.AuditLogRepository)(nil), auditRepository); err != nil {
		panic(err)
	}
}

func registerServices(container *Container) {
	userService := services.NewUserService(Resolve[repositories.UserRepository]())

	balanceService := services.NewBalanceService(
		Resolve[repositories.BalanceRepository](),
		Resolve[repositories.BalanceHistoryRepository](),
	)

	transactionService := services.NewTransactionService(
		Resolve[repositories.TransactionRepository](),
	)

	auditLogService := services.NewAuditLogService(
		Resolve[repositories.AuditLogRepository](),
	)

	if err := container.RegisterSingleton((*services.UserService)(nil), userService); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*services.BalanceService)(nil), balanceService); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*services.TransactionService)(nil), transactionService); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*services.AuditLogService)(nil), auditLogService); err != nil {
		panic(err)
	}
}

func InitDiContainer() {
	container := NewContainer()

	diContainer = container

	registerRepositories(container)
	registerServices(container)
}
