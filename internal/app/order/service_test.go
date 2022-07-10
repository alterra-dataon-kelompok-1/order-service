package order

import (
	"testing"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var orderServiceRepoInterface repository.Repository

func setupDummyOrders() []model.Order {
	oid1, _ := uuid.Parse("b49f6655-935d-48c7-8e90-2cb43e4d8a08")
	oid2, _ := uuid.Parse("09b44a8b-c0c8-4c1f-8220-ab3fe7f6b2bf")
	mid1, _ := uuid.Parse("638723dd-76d5-4e19-965c-7b541b416ccb")
	mid2, _ := uuid.Parse("07fbb648-977e-46dd-bfc8-7843fb6823da")
	mid3, _ := uuid.Parse("fae9e513-fe98-4e36-8c7b-9a1ad75ff2c5")
	uid1, _ := uuid.Parse("0f264f7e-ebb4-4010-882b-591218ac1ac8")

	orderItems1 := model.OrderItem{
		OrderID:         oid1,
		MenuID:          mid1,
		OrderItemStatus: model.Pending,
		Quantity:        1,
		Price:           6000,
	}

	orderItems2 := model.OrderItem{
		OrderID:         oid1,
		MenuID:          mid2,
		OrderItemStatus: model.Pending,
		Quantity:        1,
		Price:           6000,
	}

	orderItems3 := model.OrderItem{
		OrderID:         oid2,
		MenuID:          mid3,
		OrderItemStatus: model.Pending,
		Quantity:        2,
		Price:           6000,
	}

	return []model.Order{
		{ID: oid1, UserID: nil, OrderStatus: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 2, OrderItems: []model.OrderItem{orderItems1, orderItems2}},
		{ID: oid2, UserID: &uid1, OrderStatus: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 2, OrderItems: []model.OrderItem{orderItems3}},
	}
}

// func TestGet(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	mockOrderRepository := mock_repository.NewMockRepository(ctrl)
// 	orderService := NewService(mockOrderRepository)
//
// 	dummyOrders := setupDummyOrders()
// 	query := dto.GetRequest{}
//
// 	ctx := context.Background()
//
// 	defaultPagination := dto.PaginationInfo{}
// 	mockOrderRepository.EXPECT().GetOrders(ctx, &query).Return(&dummyOrders, &defaultPagination, nil)
//
// 	orders, err := orderService.Get(ctx, &query)
// 	log.Println("==>> orders from mock")
// 	log.Println(orders)
//
// 	require.NoError(t, err)
// 	require.Nil(t, err)
// 	require.Nil(t, orders)
//
// }

func TestSumItemQuantity(t *testing.T) {
	// Test for 0 item
	orderItemRequest := []dto.CreateOrderItemRequest{}
	res := sumItemQuantity(orderItemRequest)
	assert.Equal(t, 0, res)

	// Test for 1 item
	orderItemRequest = []dto.CreateOrderItemRequest{
		{Quantity: 10},
	}
	res = sumItemQuantity(orderItemRequest)
	assert.Equal(t, 10, res)

	// Test for more than 1 item
	orderItemRequest = []dto.CreateOrderItemRequest{
		{Quantity: 1},
		{Quantity: 2},
		{Quantity: 3},
	}
	res = sumItemQuantity(orderItemRequest)
	assert.Equal(t, 6, res)
}

func TestSumItemPrice(t *testing.T) {
	// Test for 0 item
	orderItem := []model.OrderItem{}
	res := sumItemPrice(orderItem)
	assert.Equal(t, float64(0), res)

	// Test for 1 item
	orderItem = []model.OrderItem{
		{Price: 1_000.05, Quantity: 3},
	}
	res = sumItemPrice(orderItem)
	assert.Equal(t, float64(3_000.15), res)

	// Test for more than 1 item
	orderItem = []model.OrderItem{
		{Price: 1_000.0, Quantity: 1},
		{Price: 1_000.02, Quantity: 1},
	}
	res = sumItemPrice(orderItem)
	assert.Equal(t, float64(2_000.02), res)
}
