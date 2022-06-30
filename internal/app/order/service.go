package order

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
)

type Service interface {
	Create(ctx context.Context, payload dto.CreateOrderRequest) (*model.Order, error)
}

type service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) Service {
	return &service{repository}
}

func (s *service) Create(ctx context.Context, payload dto.CreateOrderRequest) (*model.Order, error) {

	// Cannot create order if no order item in payload
	if len(payload.OrderItems) == 0 {
		return nil, errors.New("order shall have minimum 1 item")
	}

	newOrder := new(model.Order)
	newOrder.OrderItems = make([]model.OrderItem, len(payload.OrderItems))

	// assign userID
	newOrder.UserID = payload.UserID

	// assign order ID
	// TODO: implement UUID later for OrderID
	rand.Seed(time.Now().UnixNano())
	newOrder.ID = uint(rand.Intn(100))

	// new order always assigned with status 1: pending
	// TODO: add logic to implement if order made directly in cashier, it can create order with status paid
	newOrder.OrderStatusID = 1

	// Calculate Total Quantity
	newOrder.TotalQuantity = sumItemQuantity(payload.OrderItems)

	// Assign item from payload to model
	for i, item := range payload.OrderItems {
		newOrder.OrderItems[i].OrderID = &newOrder.ID
		newOrder.OrderItems[i].MenuID = item.MenuID
		newOrder.OrderItems[i].OrderItemStatusId = 1
		newOrder.OrderItems[i].Price = getItemPrice(item.MenuID)
		newOrder.OrderItems[i].Quantity = item.Quantity
	}

	// Calculate Total Price
	newOrder.TotalPrice = sumItemPrice(newOrder.OrderItems)

	createdOrder, err := s.repository.Create(ctx, *newOrder)

	return createdOrder, err
	// return newOrder, nil
}

func sumItemQuantity(s []dto.CreateOrderItemRequest) int {
	var sum int
	for _, item := range s {
		sum += item.Quantity
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

// TODO: Implement get item price (API call to other table)
func getItemPrice(id int) float32 {
	return float32(3000)
}
