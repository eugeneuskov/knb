package app

import (
	"fmt"
	"knb/app/config"
	"knb/app/handlers"
	"knb/app/repositories"
	"knb/app/services"
	"knb/db"
	"log"
)

type Application struct {
	config     *config.Config
	httpServer *httpServer
	db         *db.DB
}

func NewApplication(config *config.Config) *Application {
	return &Application{
		config:     config,
		httpServer: new(httpServer),
		db:         new(db.DB),
	}
}

func (app *Application) Run() {
	if err := app.connectDB(); err != nil {
		log.Fatal(err.Error())
	}

	if err := app.runHttpServer(); err != nil {
		log.Fatalf("Error occured while running HTTP server: %s\n", err.Error())
	}

	println("App started")
}

func (app *Application) connectDB() error {
	if err := app.db.NewPostgresDb(&app.config.DbConfig); err != nil {
		return fmt.Errorf("Failed to initialize DB: %s\n", err.Error())
	}

	if err := app.db.Migrate(); err != nil {
		return fmt.Errorf("DB Migrate is failed: %s\n", err.Error())
	}

	return nil
}

func (app *Application) runHttpServer() error {
	return app.httpServer.run(
		app.config.AppPort,
		handlers.NewHandler(
			services.NewService(
				repositories.NewRepository(app.db.DB()),
				app.config,
			),
		).InitRoutes(app.config.HandlerMode),
	)
}

func (app *Application) Shutdown() {
	println("Off")
}
