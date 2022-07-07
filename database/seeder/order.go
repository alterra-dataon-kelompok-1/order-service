package seeder

import (
	"log"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func orderTableSeeder(conn *gorm.DB) {
	// Setup uuid for referencing during API test
	oid1, _ := uuid.Parse("b49f6655-935d-48c7-8e90-2cb43e4d8a08")
	oid2, _ := uuid.Parse("09b44a8b-c0c8-4c1f-8220-ab3fe7f6b2bf")
	oid3, _ := uuid.Parse("aca1522a-07b6-4c0c-aed6-04a1d123835f")
	oid4, _ := uuid.Parse("efdd3bf4-11cc-4a3f-905b-a8e8e4c2fe72")
	oid5, _ := uuid.Parse("e085a0f9-de1f-4a01-a5d1-0cdd89fd9c9b")

	mid1, _ := uuid.Parse("b49f6655-935d-48c7-8e90-2cb43e4d8a08")
	mid2, _ := uuid.Parse("09b44a8b-c0c8-4c1f-8220-ab3fe7f6b2bf")
	mid3, _ := uuid.Parse("fae9e513-fe98-4e36-8c7b-9a1ad75ff2c5")
	mid4, _ := uuid.Parse("e6ecfbe9-975c-4aad-b6a5-30effdffb6e3")

	uid1, _ := uuid.Parse("0f264f7e-ebb4-4010-882b-591218ac1ac8")
	uid2, _ := uuid.Parse("c36c641e-4201-4257-bc5d-019d68dcd3bf")
	uid3, _ := uuid.Parse("24864159-da31-4742-90a4-eb24f01f51be")

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

	orderItems4 := model.OrderItem{
		OrderID:         oid3,
		MenuID:          mid3,
		OrderItemStatus: model.Pending,
		Quantity:        2,
		Price:           6000,
	}

	orderItems5 := model.OrderItem{
		OrderID:         oid4,
		MenuID:          mid4,
		OrderItemStatus: model.Pending,
		Quantity:        2,
		Price:           6000,
	}

	orderItems6 := model.OrderItem{
		OrderID:         oid4,
		MenuID:          mid1,
		OrderItemStatus: model.Pending,
		Quantity:        2,
		Price:           6000,
	}

	var orders = []model.Order{
		{ID: oid1, UserID: nil, OrderStatus: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 2, OrderItems: []model.OrderItem{orderItems1, orderItems2}},
		{ID: oid2, UserID: &uid1, OrderStatus: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 2, OrderItems: []model.OrderItem{orderItems3}},
		{ID: oid3, UserID: &uid2, OrderStatus: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 4, OrderItems: []model.OrderItem{orderItems4}},
		{ID: oid4, UserID: &uid2, OrderStatus: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 3, OrderItems: []model.OrderItem{orderItems5}},
		{ID: oid5, UserID: &uid3, OrderStatus: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 3, OrderItems: []model.OrderItem{orderItems6}},
	}

	if err := conn.Create(&orders).Error; err != nil {
		log.Printf("cannot seed data orders, with error %v\n", err)
	} else {
		log.Println("success seed data orders")
	}
}
