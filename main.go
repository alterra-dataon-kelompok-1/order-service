package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/alterra-dataon-kelompok-1/order-service/config"
	"github.com/alterra-dataon-kelompok-1/order-service/database"
	"github.com/alterra-dataon-kelompok-1/order-service/database/seeder"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/app/order"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/middleware"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/alterra-dataon-kelompok-1/order-service/pkg/fetcher"
	"github.com/labstack/echo/v4"
)

func main() {
	var envFlag string

	flag.StringVar(
		&envFlag,
		"env",
		"prod",
		`check if application to be run in development or production environment.
		options:
		-env=dev
		-env=prod
		`,
	)
	flag.Parse()

	config := config.New(".env")

	db, err := database.ConnectSQLDB(config)
	if err != nil {
		panic(err)
	}

	database.MigrateDB(db)
	seeder.Seed(db)

	menuFetcher := fetcher.NewFetcher(os.Getenv("MENU_SERVICE_ENDPOINT"))
	orderRepo := repository.NewRepository(db)
	orderService := order.NewService(orderRepo, menuFetcher)

	if envFlag == "dev" {
		orderService = order.NewService(orderRepo, &fetcher.MockFetcher{})
	}

	orderHandler := order.NewHandler(orderService)

	app := echo.New()
	order.RegisterHandlers(app, orderHandler)
	middleware.Init(app)

	// Set logger and close the os file
	logFile, err := os.OpenFile("logfile", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("error opening file: %v\n", err)
		middleware.SetEchoLogger(app, os.Stdout)
	} else {
		middleware.SetEchoLogger(app, logFile)
		defer logFile.Close()
	}

	// Gracefully shutdown
	go func() {
		app.Start(fmt.Sprintf(":%s", config.Get("APP_PORT")))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// During development, clean up on close
	if envFlag == "dev" {
		err = database.DropTables(db)
		if err == nil {
			log.Println("Dropping table success...")
		} else {
			log.Println("Error dropping table", err)
		}
	}

	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
