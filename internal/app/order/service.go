package order

import (
	"context"
	"errors"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/alterra-dataon-kelompok-1/order-service/pkg/utils/helper/fetcher"
	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, payload *dto.GetRequest) (*dto.SearchGetResponse[model.Order], error)
	Create(ctx context.Context, payload dto.CreateOrderRequest) (*model.Order, error)
	GetOrderByID(ctx context.Context, payload *dto.ByIDRequest) (*model.Order, error)
	DeleteOrderByID(ctx context.Context, payload *dto.ByIDRequest) (*model.Order, error)
	UpdateOrderByID(c context.Context, id uuid.UUID, payload *dto.UpdateOrderRequest) (*dto.GetOrderResponse, error)
}

type service struct {
	repository  repository.Repository
	menuFetcher fetcher.Fetcher
}

func NewService(repository repository.Repository, fetcher fetcher.Fetcher) Service {
	return &service{repository, fetcher}
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
		return nil, errors.New("E_MINIMUM_ORDER")
	}

	newOrder := new(model.Order)
	newOrder.OrderItems = make([]model.OrderItem, len(payload.OrderItems))

	// assign userID
	newOrder.UserID = payload.UserID

	// assign order ID
	newOrder.ID = uuid.New()

	// new order always assigned with status 1: pending
	// TODO: add logic to implement if order made directly in cashier, it can create order with status paid
	newOrder.OrderStatus = model.PendingOrder

	// Assign item from payload to model
	for i, item := range payload.OrderItems {
		menuDetail, err := s.menuFetcher.FetchMenuDetail(item.MenuID)
		if err != nil {
			return nil, err
		}

		if menuDetail.InStock == 0 {
			return nil, errors.New("E_NO_STOCK")
		}
		newOrder.OrderItems[i].OrderID = newOrder.ID
		newOrder.OrderItems[i].MenuID = item.MenuID
		newOrder.OrderItems[i].OrderItemStatus = model.Pending
		newOrder.OrderItems[i].Price = menuDetail.Price
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

func (s *service) DeleteOrderByID(ctx context.Context, payload *dto.ByIDRequest) (*model.Order, error) {
	data, err := s.repository.GetOrderByID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	err = s.repository.DeleteOrderByID(ctx, payload.ID)
	return data, err
}

func (s *service) UpdateOrderByID(c context.Context, id uuid.UUID, payload *dto.UpdateOrderRequest) (*dto.GetOrderResponse, error) {
	existingOrder, err := s.repository.GetOrderByID(c, id)
	if err != nil {
		return nil, err
	}

	// Prevent cancel order request when order is being prepared
	if payload.OrderStatus != nil {
		if *payload.OrderStatus == model.CanceledOrder && orderCanBeCanceled(existingOrder.OrderStatus) == false {
			return nil, errors.New("E_ORDER_CANCEL")
		}
	}

	if payload.OrderItems == nil {
		err = s.repository.UpdateOrderStatusByID(c, id, payload)
	} else {
		updateOrder := model.Order{}
		updateOrder.OrderItems = make([]model.OrderItem, len(*payload.OrderItems))
		updateOrder.ID = id
		updateOrder.UserID = existingOrder.UserID

		// Assign item from payload to model
		for i, item := range *payload.OrderItems {
			menuDetail, err := s.menuFetcher.FetchMenuDetail(item.MenuID)
			if err != nil {
				return nil, err
			}

			if menuDetail.InStock == 0 {
				return nil, errors.New("E_NO_STOCK")
			}

			updateOrder.OrderItems[i].OrderID = id
			updateOrder.OrderItems[i].MenuID = item.MenuID
			updateOrder.OrderItems[i].Price = menuDetail.Price

			if item.Quantity != nil {
				updateOrder.OrderItems[i].Quantity = *item.Quantity
			}
			if item.Status != nil {
				updateOrder.OrderItems[i].OrderItemStatus = *item.Status
			}
		}
		// Update the item
		err := s.repository.UpdateOrderByID(c, id, &updateOrder)
		if err != nil {
			return nil, err
		}
	}

	// Fetch updated order data from database
	updated, err := s.repository.GetOrderByID(c, id)
	res := newGetOrderRequestFromModel(updated)

	return &res, nil
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

func sumItemPrice(s []model.OrderItem) float64 {
	var sum float64
	for _, item := range s {
		sum = sum + (item.Price * float64(item.Quantity))
	}
	return sum
}

func orderCanBeCanceled(current model.OrderStatus) bool {
	if current == model.PendingOrder {
		return true
	}
	if current == model.PaidOrder {
		return true
	}
	return false
}

func newGetOrderRequestFromModel(m *model.Order) dto.GetOrderResponse {
	orderItems := make([]dto.GetOrderItemResponse, len(m.OrderItems))
	for i, item := range m.OrderItems {
		orderItems[i].MenuID = item.MenuID
		orderItems[i].OrderItemStatus = string(item.OrderItemStatus)
		orderItems[i].Quantity = item.Quantity
		orderItems[i].Price = item.Price
	}

	return dto.GetOrderResponse{
		ID:            m.ID,
		UserID:        m.UserID,
		OrderStatus:   m.OrderStatus,
		TotalPrice:    m.TotalPrice,
		TotalQuantity: m.TotalQuantity,
		OrderItems:    orderItems,
	}
}
