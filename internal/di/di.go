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
	if !ok || resolved == nil {
		panic(fmt.Sprintf("Failed to resolve dependency: %T", zero))
	}
	return resolved
}

func registerServices(container *Container) {
	userService := services.NewUserService()
	transactionService := services.NewTransactionService()
	balanceService := services.NewBalanceService()

	if err := container.RegisterSingleton((*services.UserService)(nil), userService); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*services.BalanceService)(nil), balanceService); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*services.TransactionService)(nil), transactionService); err != nil {
		panic(err)
	}
}

func registerRepositories(container *Container) {
	userRepository := repositories.NewUserRepository()
	balanceRepository := repositories.NewBalanceRepository()
	transactionRepository := repositories.NewTransactionRepository()

	if err := container.RegisterSingleton((*repositories.UserRepository)(nil), userRepository); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*repositories.TransactionRepository)(nil), transactionRepository); err != nil {
		panic(err)
	}

	if err := container.RegisterSingleton((*repositories.BalanceRepository)(nil), balanceRepository); err != nil {
		panic(err)
	}
}

func InitDiContainer() {
	container := NewContainer()

	registerRepositories(container)
	registerServices(container)

	diContainer = container
}
