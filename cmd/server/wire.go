//go:build wireinject
// +build wireinject

// Package main provides dependency injection setup using Wire.
package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/octop162/normal-form-app-by-claude/internal/handler"
	"github.com/octop162/normal-form-app-by-claude/internal/repository"
	"github.com/octop162/normal-form-app-by-claude/internal/service"
	"github.com/octop162/normal-form-app-by-claude/pkg/config"
	"github.com/octop162/normal-form-app-by-claude/pkg/database"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
	"github.com/octop162/normal-form-app-by-claude/pkg/validator"
)

// Provider functions for dependency injection
func provideLogger(cfg *config.Config) *logger.Logger {
	return logger.NewLogger(cfg.Log.Level)
}

func provideDB(cfg *config.Config, log *logger.Logger) (*database.DB, error) {
	return database.NewDB(&cfg.Database, log)
}

func provideSQLDB(db *database.DB) *sql.DB {
	return db.DB
}

func provideCleanupFunc(db *database.DB) func() {
	return func() {
		if db != nil {
			if err := db.Close(); err != nil {
				// Log error if possible, but don't panic during cleanup
			}
		}
	}
}

// Application holds all application components
type Application struct {
	UserHandler    *handler.UserHandler
	SessionHandler *handler.SessionHandler
	OptionHandler  *handler.OptionHandler
	AddressHandler *handler.AddressHandler
	PlanHandler    *handler.PlanHandler
	HealthHandler  *handler.HealthHandler
	DB             *sql.DB
	Logger         *logger.Logger
	Config         *config.Config
}

// Repository provider set
var repositorySet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewSessionRepository,
	repository.NewUserOptionRepository,
	repository.NewOptionRepository,
	repository.NewPrefectureRepository,
)

// Service provider set
var serviceSet = wire.NewSet(
	service.NewUserService,
	service.NewSessionService,
	service.NewOptionService,
	service.NewAddressService,
	service.NewPlanService,
)

// Handler provider set
var handlerSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewSessionHandler,
	handler.NewOptionHandler,
	handler.NewAddressHandler,
	handler.NewPlanHandler,
	handler.NewHealthHandler,
)

// Infrastructure provider set
var infrastructureSet = wire.NewSet(
	config.LoadConfig,
	provideLogger,
	provideDB,
	provideSQLDB,
	provideCleanupFunc,
	validator.NewValidator,
)

// wireApp initializes the entire application with dependency injection
func wireApp() (*Application, func(), error) {
	wire.Build(
		infrastructureSet,
		repositorySet,
		serviceSet,
		handlerSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil, nil
}