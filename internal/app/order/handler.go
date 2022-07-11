package order

import (
	"log"
	"net/http"
	"os"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	res "github.com/alterra-dataon-kelompok-1/order-service/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Get(c echo.Context) error
	Create(c echo.Context) error
	GetOrderByID(c echo.Context) error
	DeleteOrderByID(c echo.Context) error
	UpdateOrderByID(c echo.Context) error
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

	orders, err := h.service.Get(c.Request().Context(), &payload)
	if err != nil {
		return res.NewErrorResponse(c, res.ErrorConst.InternalServerError)
	}
	return c.JSON(http.StatusOK, orders)
}

func (h *handler) Create(c echo.Context) error {
	payload := new(dto.CreateOrderRequest)
	if err := c.Bind(&payload); err != nil {
		log.Println("err:", err)
		return res.NewErrorResponse(c, res.ErrorConst.BadRequest)
	}
	log.Println("payload")

	newOrder, err := h.service.Create(c.Request().Context(), *payload)
	if err != nil {
		log.Println("err:", err)
		return res.NewErrorResponse(c, res.ErrorConst.BadRequest)
	}
	log.Println("newOrder")

	// TODO: Implement passing resource location during runtime
	successResp := res.NewSuccessBuilder().Status(http.StatusCreated).WithData(*newOrder).WithResourceLocation(os.Getenv("APP_HOST_URL"), newOrder.ID.String())
	return successResp.SendJSON(c)
}

func (h *handler) GetOrderByID(c echo.Context) error {
	payload := dto.ByIDRequest{}
	if err := c.Bind(&payload); err != nil {
		return res.NewErrorResponse(c, res.ErrorConst.BadRequest)
	}

	result, err := h.service.GetOrderByID(c.Request().Context(), &payload)
	if err != nil {
		return res.NewErrorResponse(c, res.ErrorConst.ResourceNotFound)
	}

	resp := res.NewSuccessBuilder().Status(http.StatusOK).WithData(result)
	return resp.SendJSON(c)
}

func (h *handler) DeleteOrderByID(c echo.Context) error {
	payload := dto.ByIDRequest{}
	if err := c.Bind(&payload); err != nil {
		return res.NewErrorResponse(c, res.ErrorConst.BadRequest)
	}

	if err := c.Validate(payload); err != nil {
		return res.NewErrorResponse(c, res.ErrorConst.Validation)
	}

	data, err := h.service.DeleteOrderByID(c.Request().Context(), &payload)
	if err != nil {
		return res.NewErrorResponse(c, res.ErrorConst.ResourceNotFound)
	}

	resp := res.NewSuccessBuilder().Status(http.StatusOK).WithData(data)
	return resp.SendJSON(c)
}

func (h *handler) UpdateOrderByID(c echo.Context) error {
	payload := new(dto.UpdateOrderRequest)
	err := c.Bind(&payload)
	if err != nil {
		log.Println("==> error binding", err)
		return res.NewErrorResponse(c, res.ErrorConst.BadRequest)
	}
	log.Println("==> binded payload", payload)

	// err = c.Validate(&payload)
	// if err != nil {
	// 	log.Println("==> error validating", err)
	// 	return res.NewErrorResponse(c, res.ErrorConst.Validation)
	// }

	strID := c.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		log.Println("==> error:", err)
		return res.NewErrorResponse(c, res.ErrorConst.BadRequest)
	}

	data, err := h.service.UpdateOrderByID(c.Request().Context(), id, payload)
	if err != nil {
		log.Println("==> error:", err)
		if err.Error() == "E_NOT_FOUND" {
			return res.NewErrorResponse(c, res.ErrorConst.ResourceNotFound)
		}
		return res.NewErrorResponse(c, res.ErrorConst.InternalServerError)
	}

	resp := res.NewSuccessBuilder().Status(http.StatusOK).WithData(data)
	return resp.SendJSON(c)
}
