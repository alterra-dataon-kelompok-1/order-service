package main

import (
	"fmt"

	"github.com/alterra-dataon-kelompok-1/order-service/config"
	"github.com/alterra-dataon-kelompok-1/order-service/database"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/app/order"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	config := config.New(".env")

	db, err := database.ConnectSQLDB(config)
	if err != nil {
		panic(err)
	}

	database.MigrateDB(db)

	orderRepo := repository.NewRepository(db)
	orderService := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderService)

	app := echo.New()
	order.RegisterHandlers(app, orderHandler)

	app.Start(fmt.Sprintf(":%s", config.Get("APP_PORT")))
}
