// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/AlmazDefourten/goapp/container"
	"github.com/AlmazDefourten/goapp/handlers"
	"github.com/AlmazDefourten/goapp/models"
	"github.com/AlmazDefourten/goapp/services"
)

// Injectors from wire.go:

// Initialize container with global app dependencies -
// Connection, configurator, etc...
func InitializeContainer() models.Container {
	viper := container.NewViperConfigurator()
	db := container.NewConnection(viper)
	modelsContainer := container.NewContainer(db, viper)
	return modelsContainer
}

// Initialize dependencies for services
func InitServiceDependency(container_inited *models.Container) models.ServiceContainer {
	userService := services.NewUserService(container_inited)
	serviceContainer := container.NewServiceContainer(userService)
	return serviceContainer
}

// Initialize dependencies for handlers
func InitHandlerDependency(userService models.IUserService) container.HandlerContainer {
	userInfoHandler := handlers.NewUserInfoHandler(userService)
	handlerContainer := container.NewHandlerContainer(userInfoHandler)
	return handlerContainer
}

// wire.go:

// RegisterServices - decomposition ServiceContainer to services
func RegisterServices(serviceContainer *models.ServiceContainer) container.HandlerContainer {
	return InitHandlerDependency(
		serviceContainer.UserService,
	)
}