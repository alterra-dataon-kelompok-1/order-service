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
	"github.com/alterra-dataon-kelompok-1/order-service/pkg/utils/helper/fetcher"
	"github.com/google/uuid"
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
	orderService := NewService(orderRepo, &fetcher.MockFetcher{})
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

	// id declared due to request need a pointer data
	id := uuid.New()

	testCase := []struct {
		Case             string
		RequestBody      dto.CreateOrderRequest
		WantResponseCode int
		WantBodyContains string
	}{
		{
			Case: "success",
			RequestBody: dto.CreateOrderRequest{
				UserID: &id,
				OrderItems: []dto.CreateOrderItemRequest{
					{
						MenuID:   uuid.New(),
						Quantity: 1,
					},
				},
			},
			WantResponseCode: http.StatusCreated,
			WantBodyContains: id.String(),
		},
		{
			Case: "no order item",
			RequestBody: dto.CreateOrderRequest{
				UserID: &id,
				OrderItems: []dto.CreateOrderItemRequest{
					{
						MenuID:   uuid.New(),
						Quantity: 0,
					},
				},
			},
			WantResponseCode: http.StatusBadRequest,
			WantBodyContains: " request",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.Case, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.RequestBody)
			req := httptest.NewRequest("POST", "/orders", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Assertion
			if assert.NoError(t, h.Create(c)) {
				assert.Equal(t, tc.WantResponseCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.WantBodyContains)
			}
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	testCase := []struct {
		Case             string
		RequestBody      dto.UpdateOrderRequest
		ParamID          string
		WantResponseCode int
		WantBodyContains string
	}{
		{
			Case:             "success",
			ParamID:          "aca1522a-07b6-4c0c-aed6-04a1d123835f",
			WantResponseCode: http.StatusOK,
			WantBodyContains: "aca1522a-07b6-4c0c-aed6-04a1d123835f",
		},
		{
			Case:             "incorrect UUID format",
			ParamID:          "aca12a-07b6-4c0c-aed6-0423835f",
			WantResponseCode: http.StatusBadRequest,
			WantBodyContains: "E_BAD_REQUEST",
		},
		{
			Case:             "ID not found",
			ParamID:          "1c08b996-92bb-4c09-aa3b-989b4c5092ca",
			WantResponseCode: http.StatusNotFound,
			WantBodyContains: "E_NOT_FOUND",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.Case, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/orders", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ParamID)

			// Assertion
			if assert.NoError(t, h.GetOrderByID(c)) {
				assert.Equal(t, tc.WantResponseCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.WantBodyContains)
			}
		})
	}
}

func TestDeleteOrderByID_Base(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	testCase := []struct {
		Case                   string
		ParamID                string
		WantResponseCode       int
		WantAfterDeleteRequest bool
		WantBodyContains       string
	}{
		{
			Case:                   "success",
			ParamID:                "aca1522a-07b6-4c0c-aed6-04a1d123835f",
			WantResponseCode:       http.StatusOK,
			WantAfterDeleteRequest: false,
			WantBodyContains:       "aca1522a-07b6-4c0c-aed6-04a1d123835f",
		},
		{
			Case:                   "incorrect ID format",
			ParamID:                "c1522a-07b6-4c0c-aed6-04a1d1238f",
			WantResponseCode:       http.StatusBadRequest,
			WantAfterDeleteRequest: false,
			WantBodyContains:       "bad request",
		},
		{
			Case:                   "not found",
			ParamID:                "1c08b996-92bb-4c09-aa3b-989b4c5092ca",
			WantResponseCode:       http.StatusNotFound,
			WantAfterDeleteRequest: false,
			WantBodyContains:       "not found",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.Case, func(t *testing.T) {
			// t.Parallel()
			req := httptest.NewRequest("DELETE", "/orders", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ParamID)

			// Assertion
			if assert.NoError(t, h.DeleteOrderByID(c)) {
				assert.Equal(t, tc.WantResponseCode, rec.Code)
				assert.Equal(t, tc.WantAfterDeleteRequest, FindAfterDelete(e, h, tc.ParamID))
				assert.Contains(t, rec.Body.String(), tc.WantBodyContains)
			}
		})
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

func TestUpdateOrderStatus(t *testing.T) {
	// Setup
	e, db, h := createTestApp()
	defer database.DropTables(db)

	// paidOrder is declared since request body require a pointer
	paidOrder := model.PaidOrder

	testCase := []struct {
		Case             string
		RequestBody      dto.UpdateOrderRequest
		ParamID          string
		WantResponseCode int
		WantBodyContains string
	}{
		{
			Case:    "success",
			ParamID: "aca1522a-07b6-4c0c-aed6-04a1d123835f",
			RequestBody: dto.UpdateOrderRequest{
				OrderStatus: &paidOrder,
			},
			WantResponseCode: http.StatusOK,
			WantBodyContains: "paid",
		},
		{
			Case:    "not found",
			ParamID: "1c08b996-92bb-4c09-aa3b-989b4c5092ca",
			RequestBody: dto.UpdateOrderRequest{
				OrderStatus: &paidOrder,
			},
			WantResponseCode: http.StatusNotFound,
			WantBodyContains: "not found",
		},
		{
			Case:    "bad uuid format",
			ParamID: "10b996-92bb-4c09-aa3b-989b4c5092ca",
			RequestBody: dto.UpdateOrderRequest{
				OrderStatus: &paidOrder,
			},
			WantResponseCode: http.StatusBadRequest,
			WantBodyContains: "E_BAD_REQUEST",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.Case, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.RequestBody)
			req := httptest.NewRequest("PUT", "/orders", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ParamID)

			// Assertion
			if assert.NoError(t, h.UpdateOrderByID(c)) {
				assert.Equal(t, tc.WantResponseCode, rec.Code)
				assert.Contains(t, rec.Body.String(), tc.WantBodyContains)
			}
		})
	}
}
