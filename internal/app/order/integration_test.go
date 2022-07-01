package order

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alterra-dataon-kelompok-1/order-service/database"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func createTestApp() (*echo.Echo, Handler) {
	// TODO: use mock / test db instead of actual db
	// Since we do not use any mock in this test, it is more of integration test that cover service and also repository at the same time
	db, err := database.ConnectDevDB()

	if err != nil {
		panic(err)
	}

	database.MigrateDB(db)

	e := echo.New()
	orderRepo := repository.NewRepository(db)
	orderService := NewService(orderRepo)
	orderHandler := NewHandler(orderService)

	return e, orderHandler
}

func TestGetBooksHandler(t *testing.T) {
	// Setup
	e, h := createTestApp()

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertion
	if assert.NoError(t, h.Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Code)
	}
}
