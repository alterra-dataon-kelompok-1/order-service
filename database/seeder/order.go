package seeder

import (
	"log"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func orderTableSeeder(conn *gorm.DB) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	id4 := uuid.New()
	id5, _ := uuid.Parse("b8a36547-d74d-4186-b293-9aae9f87f4f3")

	orderItems1 := model.OrderItem{
		OrderID:  id1,
		MenuID:   2,
		Status:   model.Pending,
		Quantity: 1,
		Price:    6000,
	}

	orderItems2 := model.OrderItem{
		OrderID:  id1,
		MenuID:   3,
		Status:   model.Pending,
		Quantity: 1,
		Price:    6000,
	}

	orderItems3 := model.OrderItem{
		OrderID:  id2,
		MenuID:   3,
		Status:   model.Pending,
		Quantity: 2,
		Price:    6000,
	}

	orderItems4 := model.OrderItem{
		OrderID:  id3,
		MenuID:   4,
		Status:   model.Pending,
		Quantity: 2,
		Price:    6000,
	}

	orderItems5 := model.OrderItem{
		OrderID:  id4,
		MenuID:   4,
		Status:   model.Pending,
		Quantity: 2,
		Price:    6000,
	}

	orderItems6 := model.OrderItem{
		OrderID:  id4,
		MenuID:   4,
		Status:   model.Pending,
		Quantity: 2,
		Price:    6000,
	}

	var orders = []model.Order{
		{ID: id1, UserID: 0, Status: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 2, OrderItems: []model.OrderItem{orderItems1, orderItems2}},
		{ID: id2, UserID: 123, Status: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 2, OrderItems: []model.OrderItem{orderItems3}},
		{ID: id3, UserID: 11, Status: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 4, OrderItems: []model.OrderItem{orderItems4}},
		{ID: id4, UserID: 54, Status: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 3, OrderItems: []model.OrderItem{orderItems5}},
		{ID: id5, UserID: 4, Status: model.PendingOrder, TotalPrice: 12_000, TotalQuantity: 3, OrderItems: []model.OrderItem{orderItems6}},
	}

	if err := conn.Create(&orders).Error; err != nil {
		log.Printf("cannot seed data orders, with error %v\n", err)
	}
	log.Println("success seed data orders")
}
