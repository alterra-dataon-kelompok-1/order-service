package order

import (
	"context"
	"errors"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/alterra-dataon-kelompok-1/order-service/pkg/utils/helper"
	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, payload *dto.GetRequest) (*dto.SearchGetResponse[model.Order], error)
	Create(ctx context.Context, payload dto.CreateOrderRequest) (*model.Order, error)
	GetOrderByID(ctx context.Context, payload *dto.ByIDRequest) (*model.Order, error)
}

type service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) Service {
	return &service{repository}
}

func (s *service) Get(ctx context.Context, payload *dto.GetRequest) (*dto.SearchGetResponse[model.Order], error) {
	orders, paginationInfo, err := s.repository.GetOrders(ctx, payload)
	if err != nil {
		return nil, err
	}

	result := dto.SearchGetResponse[model.Order]{
		PaginationInfo: paginationInfo,
		Data:           *orders,
	}

	return &result, nil
}

func (s *service) Create(ctx context.Context, payload dto.CreateOrderRequest) (*model.Order, error) {

	// Cannot create order if no order item in payload
	reqOrderQuantity := sumItemQuantity(payload.OrderItems)
	if reqOrderQuantity == 0 {
		return nil, errors.New("order shall have minimum 1 item")
	}

	newOrder := new(model.Order)
	newOrder.OrderItems = make([]model.OrderItem, len(payload.OrderItems))

	// assign userID
	newOrder.UserID = payload.UserID

	// assign order ID
	newOrder.ID = uuid.New()

	// new order always assigned with status 1: pending
	// TODO: add logic to implement if order made directly in cashier, it can create order with status paid
	newOrder.Status = model.PendingOrder

	// Assign item from payload to model
	for i, item := range payload.OrderItems {
		newOrder.OrderItems[i].OrderID = newOrder.ID
		newOrder.OrderItems[i].MenuID = item.MenuID
		newOrder.OrderItems[i].Status = model.Pending
		newOrder.OrderItems[i].Price = helper.GetItemPrice(item.MenuID)
		newOrder.OrderItems[i].Quantity = item.Quantity
	}

	// Calculate Total Quantity and Total Price
	newOrder.TotalQuantity = sumItemQuantity(payload.OrderItems)
	newOrder.TotalPrice = sumItemPrice(newOrder.OrderItems)

	// Create order record in repository
	createdOrder, err := s.repository.Create(ctx, *newOrder)

	return createdOrder, err
}

func (s *service) GetOrderByID(ctx context.Context, payload *dto.ByIDRequest) (*model.Order, error) {
	data, err := s.repository.GetOrderByID(ctx, payload.ID)
	if err != nil {
		if err == errors.New("E_NOT_FOUND") {
			return nil, err
		}
		return nil, errors.New("E_SERVER")
	}

	// TODO: decide if we need to transfer response dto instead
	return data, nil
}

type hasQuantity interface {
	GetQuantity() int
}

func sumItemQuantity[T hasQuantity](s []T) int {
	var sum int
	for _, item := range s {
		sum += item.GetQuantity()
	}
	return sum
}

func sumItemPrice(s []model.OrderItem) float32 {
	var sum float32
	for _, item := range s {
		sum = sum + (item.Price * float32(item.Quantity))
	}
	return sum
}
