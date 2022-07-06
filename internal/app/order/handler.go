package order

import (
	"log"
	"net/http"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	res "github.com/alterra-dataon-kelompok-1/order-service/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Get(c echo.Context) error
	Create(c echo.Context) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{service}
}

// TODO: Implement sorting by field
// TODO: Implement filter by day
func (h *handler) Get(c echo.Context) error {
	var payload dto.GetRequest
	c.Bind(&payload)

	log.Println("><><><> payload")
	log.Println(payload)

	orders, err := h.service.Get(c.Request().Context(), &payload)
	if err != nil {
		return res.NewErrorResponse(c, res.ErrorConst.InternalServerError)
	}
	return c.JSON(http.StatusOK, orders)
}

func (h *handler) Create(c echo.Context) error {
	payload := new(dto.CreateOrderRequest)
	c.Bind(&payload)

	newOrder, err := h.service.Create(c.Request().Context(), *payload)
	if err != nil {
		return res.NewErrorResponse(c, res.ErrorConst.BadRequest)
	}

	// TODO: Implement passing resource location during runtime
	successResp := res.NewSuccessBuilder().Status(http.StatusCreated).WithData(*newOrder).WithResourceLocation("localhost:8050", newOrder.ID.String())
	return successResp.SendJSON(c)
}
