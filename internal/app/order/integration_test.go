package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alterra-dataon-kelompok-1/order-service/config"
	"github.com/alterra-dataon-kelompok-1/order-service/database"
	"github.com/alterra-dataon-kelompok-1/order-service/database/seeder"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/middleware"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/alterra-dataon-kelompok-1/order-service/pkg/response"
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

	middleware.Init(e)

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

func TestGetOrderByID_Base(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	req := httptest.NewRequest("GET", "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	// Following orderID is coming from seeder
	orderID := "b8a36547-d74d-4186-b293-9aae9f87f4f3"

	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(orderID)

	// Assertion
	if assert.NoError(t, h.GetOrderByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.String())

		// interface{} map[string]interface{} => model.Order

		jsonRes := response.SuccessResponse{Data: model.Order{}}
		err := json.NewDecoder(rec.Body).Decode(&jsonRes)
		if err != nil {
			t.Error(err)
		}

		// data := jsonRes.Data.(model.Order)
	}
}

func TestGetOrderByID_IncorrectUUIDFormat(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	req := httptest.NewRequest("GET", "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	// UUID is shorter than it should
	orderID := "b8a36547-d74d-4186-b293-9aae"

	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(orderID)

	// Assertion
	if assert.NoError(t, h.GetOrderByID(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	}
}

func TestGetOrderByID_NotFound(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	req := httptest.NewRequest("GET", "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	// Random UUID as input
	orderID := "1c08b996-92bb-4c09-aa3b-989b4c5092ca"

	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(orderID)

	// Assertion
	if assert.NoError(t, h.GetOrderByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	}
}

func TestDeleteOrderByID_Base(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	req := httptest.NewRequest("GET", "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	// Following orderID is coming from seeder
	orderID := "b8a36547-d74d-4186-b293-9aae9f87f4f3"

	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(orderID)

	// Assertion
	if assert.NoError(t, h.DeleteOrderByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.NotEmpty(t, rec.Body.String())
		assert.Equal(t, false, FindAfterDelete(e, h, orderID))
	}
}

func TestDeleteOrderByID_WrongInputFormat(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	req := httptest.NewRequest("DELETE", "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	// Following orderID is coming from seeder
	orderID := "b8a36547-d74d-4186-b293-9aae"

	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(orderID)

	// Assertion
	if assert.NoError(t, h.DeleteOrderByID(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestDeleteOrderByID_NotFound(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	req := httptest.NewRequest("DELETE", "/orders", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	// Following orderID is coming from seeder
	orderID := "1c08b996-92bb-4c09-aa3b-989b4c5092ca"

	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(orderID)

	// Assertion
	if assert.NoError(t, h.DeleteOrderByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

// FindAfterDelete is a helper function to check
func FindAfterDelete(e *echo.Echo, h Handler, stringUUID string) bool {
	req := httptest.NewRequest("GET", "/orders", nil)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(stringUUID)

	h.GetOrderByID(c)
	fmt.Println("status:", rec.Code)
	if rec.Code == http.StatusOK {
		return true
	}
	return false
}
