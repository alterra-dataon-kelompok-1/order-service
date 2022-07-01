package seeder

import (
	"log"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"gorm.io/gorm"
)

func orderStatusesTableSeeder(conn *gorm.DB) {
	var orderStatuses = []model.OrderStatus{
		{ID: 1, State: "pending"},
		{ID: 2, State: "paid"},
		{ID: 3, State: "prepared"},
		{ID: 4, State: "completed"},
		{ID: 5, State: "canceled"},
	}

	if err := conn.Table("order_statuses").Create(&orderStatuses).Error; err != nil {
		log.Printf("cannot seed table order_statuses, with error %v\n", err)
	}
	log.Println("success seed table order_statuses")
}
