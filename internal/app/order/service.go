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
	DeleteOrderByID(ctx context.Context, payload *dto.ByIDRequest) (*model.Order, error)
	UpdateOrderByID(c context.Context, id uuid.UUID, payload *dto.UpdateOrderRequest) error
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
	newOrder.OrderStatus = model.PendingOrder

	// Assign item from payload to model
	for i, item := range payload.OrderItems {
		newOrder.OrderItems[i].OrderID = newOrder.ID
		newOrder.OrderItems[i].MenuID = item.MenuID
		newOrder.OrderItems[i].OrderItemStatus = model.Pending
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

func (s *service) DeleteOrderByID(ctx context.Context, payload *dto.ByIDRequest) (*model.Order, error) {
	data, err := s.repository.GetOrderByID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	err = s.repository.DeleteOrderByID(ctx, payload.ID)
	return data, err
}

func (s *service) UpdateOrderByID(c context.Context, id uuid.UUID, payload *dto.UpdateOrderRequest) error {
	// TODO: add logic to prevent cancel order when order is being prepared
	// queriedOrder, err := s.repository.GetOrderByID(c, id)
	// if err != nil {
	// 	return err
	// }

	// FIX: issue might be due to golang not allow us to
	// fmt.Println(*payload.OrderStatus == model.CanceledOrder)
	// fmt.Println(queriedOrder.OrderStatus != model.PendingOrder)
	// fmt.Println(queriedOrder.OrderStatus != model.PaidOrder)
	// if *payload.OrderStatus == model.CanceledOrder {
	// 	if oldData.OrderStatus != model.PendingOrder || oldData.OrderStatus != model.PaidOrder {
	// 		return errors.New("cannot cancel order after prepared")
	// 	}
	// }

	/* orderData := make(map[string]interface{})
	if payload.OrderStatus != nil {
		orderData["order_status"] = payload.OrderStatus

		err := s.repository.UpdateOrderByID(c, id, orderData)
		if err != nil {
			return errors.New("E_SERVER")
		}
	} */

	var update model.Order
	if payload.OrderStatus != nil {
		update.OrderStatus = *payload.OrderStatus
	}
	if payload.OrderItems != nil {
		update.OrderItems = make([]model.OrderItem, len(*payload.OrderItems))
		for i, item := range *payload.OrderItems {
			update.OrderItems[i].OrderID = id
			update.OrderItems[i].MenuID = item.MenuID
			if item.Status != nil {
				update.OrderItems[i].OrderItemStatus = *item.Status
			}
			if item.Quantity != nil {
				update.OrderItems[i].Quantity = *item.Quantity
			}
		}
	}

	err := s.repository.UpdateOrderByIDWithModel(c, id, &update)
	if err != nil {
		return err
	}

	// if payload.OrderItems != nil {
	// 	for _, v := range *payload.OrderItems {
	// 		log.Println(v)
	// 		orderItemData := make(map[string]interface{})
	// 		if v.Status != nil {
	// 			orderItemData["order_item_status"] = v.Status
	// 		}
	// 		if v.Quantity != nil {
	// 			orderItemData["quantity"] = v.Quantity
	// 		}
	// 		fmt.Println("orderitems from payload=>", orderItemData)
	// 	}
	// }

	return nil
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
