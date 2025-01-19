package di

import (
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

func Resolve(serviceType any) any {
	t := reflect.TypeOf(serviceType).Elem()

	if resolved, exists := diContainer.singleton[t]; exists {
		return resolved.Interface()
	}

	return nil
}

func InitDiContainer() {
	container := NewContainer()

	userService := services.NewUserService()

	if err := container.RegisterSingleton((*services.UserService)(nil), userService); err != nil {
		panic(err)
	}

	diContainer = container
}
