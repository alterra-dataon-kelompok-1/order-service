package order

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alterra-dataon-kelompok-1/order-service/config"
	"github.com/alterra-dataon-kelompok-1/order-service/database"
	"github.com/alterra-dataon-kelompok-1/order-service/database/seeder"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createTestApp() (*echo.Echo, *gorm.DB, Handler) {
	// TODO: use mock / test db instead of actual db
	// Since we do not use any mock in this test, it is more of integration test that cover service and also repository at the same time
	// db, err := database.ConnectDevDB()
	db, err := database.ConnectSQLDB(config.New("../../../.env"))

	if err != nil {
		panic(err)
	}

	database.MigrateDB(db)
	seeder.Seed(db)

	e := echo.New()
	orderRepo := repository.NewRepository(db)
	orderService := NewService(orderRepo)
	orderHandler := NewHandler(orderService)

	return e, db, orderHandler
}

func TestGetOrders_Base(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	req := httptest.NewRequest("GET", "/orders", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// Assertion
	if assert.NoError(t, h.Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.String())

		jsonRes := dto.SearchGetResponse[model.Order]{}
		err := json.NewDecoder(rec.Body).Decode(&jsonRes)
		if err != nil {
			t.Error(err)
		}
		assert.NotZero(t, len(jsonRes.Data))
		assert.Equal(t, len(jsonRes.Data), jsonRes.PaginationInfo.Count)
	}
}

func TestGetOrders_WithPagination(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	req := httptest.NewRequest("GET", "/orders?page_size=2", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// Assertion
	if assert.NoError(t, h.Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.String())

		res := dto.SearchGetResponse[model.Order]{PaginationInfo: nil}
		err := json.NewDecoder(rec.Body).Decode(&res)
		if err != nil {
			t.Error(err)
		}

		assert.NotZero(t, len(res.Data))
		assert.Equal(t, 2, *res.PaginationInfo.PageSize)
	}
}

func TestCreateOrder(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	newOrder := dto.CreateOrderRequest{
		UserID: 999,
		OrderItems: []dto.CreateOrderItemRequest{
			{
				MenuID:   123,
				Quantity: 1,
			},
		},
	}
	reqBody, _ := json.Marshal(newOrder)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/orders")

	// Assertion
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

// Tested Endpoint: POST /orders/:id
// Payload Constrain: ID not provided
// Expected Behaviour:
// - Code: 400
func TestCreateOrder_NoOrderItem(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	newOrder := dto.CreateOrderRequest{
		UserID: 999,
		OrderItems: []dto.CreateOrderItemRequest{
			{
				MenuID:   123,
				Quantity: 0,
			},
		},
	}
	reqBody, _ := json.Marshal(newOrder)

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// Assertion
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}
