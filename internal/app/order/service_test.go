package order

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mock_repository "github.com/alterra-dataon-kelompok-1/order-service/internal/app/order/mock"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/repository"
	"github.com/alterra-dataon-kelompok-1/order-service/pkg/fetcher"
	"github.com/golang/mock/gomock"
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

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dummyOrders := setupDummyOrders()
	defaultPage := 1
	defaultPageSize := 10

	type args struct {
		ctx     context.Context
		payload *dto.GetRequest
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(orderRepo *mock_repository.MockRepository)
		want       *dto.SearchGetResponse[model.Order]
		wantErr    bool
	}{
		{
			name: "success get all user w/o pagination",
			args: args{
				ctx:     context.TODO(),
				payload: &dto.GetRequest{},
			},
			beforeTest: func(orderRepo *mock_repository.MockRepository) {
				orderRepo.EXPECT().GetOrders(context.TODO(), &dto.GetRequest{}).Return(
					&dummyOrders,
					&dto.PaginationInfo{
						Pagination: dto.Pagination{
							Page:     &defaultPage,
							PageSize: &defaultPageSize,
						},
						Count:       2,
						MoreRecords: false,
						TotalPage:   1,
					},
					nil,
				)
			},
			want: &dto.SearchGetResponse[model.Order]{
				Data: dummyOrders,
				PaginationInfo: &dto.PaginationInfo{Pagination: dto.Pagination{
					Page:     &defaultPage,
					PageSize: &defaultPageSize,
				}, Count: 2, MoreRecords: false, TotalPage: 1},
			},
			wantErr: false,
		},
		{
			name: "error from repo",
			args: args{
				ctx:     context.TODO(),
				payload: &dto.GetRequest{},
			},
			beforeTest: func(orderRepo *mock_repository.MockRepository) {
				orderRepo.EXPECT().GetOrders(context.TODO(), &dto.GetRequest{}).Return(nil, nil, errors.New("simulated db err"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := mock_repository.NewMockRepository(ctrl)
			orderService := NewService(mockOrderRepo, &fetcher.MockFetcher{})

			if tt.beforeTest != nil {
				tt.beforeTest(mockOrderRepo)
			}

			got, err := orderService.Get(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uid := uuid.New()

	type args struct {
		ctx     context.Context
		payload dto.CreateOrderRequest
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(orderRepo *mock_repository.MockRepository)
		want       *model.Order
		wantErr    bool
	}{
		{
			name: "no order item provided",
			args: args{
				ctx: context.TODO(),
				payload: dto.CreateOrderRequest{
					UserID:     &uid,
					OrderItems: nil,
				},
			},
			beforeTest: func(orderRepo *mock_repository.MockRepository) {
				orderRepo.EXPECT().Create(context.TODO(), &model.Order{}).Times(0)
			},
			want:    nil,
			wantErr: true,
		},
		// NOTE: success cannot be unit tested due to random UUID is generated in the function
		// {
		// 	name: "success",
		// 	args: args{
		// 		ctx: context.TODO(),
		// 		payload: dto.CreateOrderRequest{
		// 			UserID: &uid,
		// 			OrderItems: []dto.CreateOrderItemRequest{
		// 				{
		// 					MenuID:   mid1,
		// 					Quantity: 1,
		// 				},
		// 				{
		// 					MenuID:   mid2,
		// 					Quantity: 2,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	beforeTest: func(orderRepo *mock_repository.MockRepository) {
		// 		orderRepo.EXPECT().Create(context.TODO(), &model.Order{
		// 			ID:            uuid.UUID{},
		// 			UserID:        &uid,
		// 			OrderStatus:   model.PendingOrder,
		// 			TotalPrice:    30_000,
		// 			TotalQuantity: 3,
		// 			OrderItems: []model.OrderItem{
		// 				{
		// 					OrderID:         uuid.UUID{},
		// 					MenuID:          mid1,
		// 					OrderItemStatus: model.Pending,
		// 					Quantity:        1,
		// 					Price:           10_000,
		// 					Model:           model.Model{},
		// 				},
		// 				{
		// 					OrderID:         uuid.UUID{},
		// 					MenuID:          mid2,
		// 					OrderItemStatus: model.Pending,
		// 					Quantity:        2,
		// 					Price:           10_000,
		// 					Model:           model.Model{},
		// 				},
		// 			},
		// 			Model: model.Model{},
		// 		}).Return(&model.Order{
		// 			ID:            uuid.UUID{},
		// 			UserID:        &uid,
		// 			OrderStatus:   model.PendingOrder,
		// 			TotalPrice:    30_000,
		// 			TotalQuantity: 3,
		// 			OrderItems: []model.OrderItem{
		// 				{
		// 					OrderID:         uuid.UUID{},
		// 					MenuID:          mid1,
		// 					OrderItemStatus: model.Pending,
		// 					Quantity:        1,
		// 					Price:           10_000,
		// 					Model:           model.Model{},
		// 				},
		// 				{
		// 					OrderID:         uuid.UUID{},
		// 					MenuID:          mid2,
		// 					OrderItemStatus: model.Pending,
		// 					Quantity:        2,
		// 					Price:           10_000,
		// 					Model:           model.Model{},
		// 				},
		// 			},
		// 			Model: model.Model{},
		// 		},
		// 			nil)
		// 	},
		// 	want: &model.Order{
		// 		// ID:            oid,
		// 		UserID:        &uid,
		// 		OrderStatus:   model.PendingOrder,
		// 		TotalPrice:    30_000,
		// 		TotalQuantity: 3,
		// 		OrderItems: []model.OrderItem{
		// 			{
		// 				// OrderID:         oid,
		// 				MenuID:          mid1,
		// 				OrderItemStatus: model.Pending,
		// 				Quantity:        1,
		// 				Price:           10_000,
		// 				Model:           model.Model{},
		// 			},
		// 			{
		// 				// OrderID:         oid,
		// 				MenuID:          mid2,
		// 				OrderItemStatus: model.Pending,
		// 				Quantity:        2,
		// 				Price:           10_000,
		// 				Model:           model.Model{},
		// 			},
		// 		},
		// 		Model: model.Model{},
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := mock_repository.NewMockRepository(ctrl)
			orderService := NewService(mockOrderRepo, &fetcher.MockFetcher{})

			if tt.beforeTest != nil {
				tt.beforeTest(mockOrderRepo)
			}

			got, err := orderService.Create(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestOrderCanBeCanceled(t *testing.T) {
	assert.Equal(t, true, orderCanBeCanceled(model.PendingOrder))
	assert.Equal(t, true, orderCanBeCanceled(model.PaidOrder))
	assert.Equal(t, false, orderCanBeCanceled(model.CanceledOrder))
}

func TestModelToGetOrderResponse(t *testing.T) {
	orderItems5 := model.OrderItem{
		OrderID:         uuid.New(),
		MenuID:          uuid.New(),
		OrderItemStatus: model.Pending,
		Quantity:        2,
		Price:           6000,
	}

	orderItems6 := model.OrderItem{
		OrderID:         uuid.New(),
		MenuID:          uuid.New(),
		OrderItemStatus: model.Pending,
		Quantity:        2,
		Price:           6000,
	}

	uid := uuid.New()
	o := model.Order{
		ID:            uuid.New(),
		UserID:        &uid,
		OrderStatus:   model.CanceledOrder,
		TotalPrice:    123.4,
		TotalQuantity: 3,
		OrderItems: []model.OrderItem{
			orderItems5,
			orderItems6,
		},
	}

	r := ModelToGetOrderResponse(&o)

	assert.Equal(t, o.ID, r.ID)
	assert.Equal(t, o.UserID, r.UserID)
	assert.Equal(t, o.OrderStatus, r.OrderStatus)
	assert.Equal(t, o.TotalPrice, r.TotalPrice)
	assert.Equal(t, o.TotalQuantity, r.TotalQuantity)
}
