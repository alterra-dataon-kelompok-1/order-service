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
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/alterra-dataon-kelompok-1/order-service/pkg/utils/interservice"
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

	if envFlag == "dev" {
		seeder.Seed(db)
	}

	orderRepo := repository.NewRepository(db)
	orderService := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderService)

	// FIX: how to implement cleanly
	menuServiceAPI := interservice.NewInterservice(config)
	log.Println(menuServiceAPI)

	app := echo.New()
	order.RegisterHandlers(app, orderHandler)

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
