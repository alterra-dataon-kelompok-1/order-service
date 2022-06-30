package seeder

import (
	"log"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"gorm.io/gorm"
)

func orderItemStatusesTableSeeder(conn *gorm.DB) {
	var orderItemStatuses = []model.OrderItemStatus{
		{ID: 1, State: "pending"},
		{ID: 2, State: "paid"},
		{ID: 3, State: "prepared"},
		{ID: 4, State: "delivered"},
	}

	if err := conn.Table("order_item_statuses").Create(&orderItemStatuses).Error; err != nil {
		log.Printf("cannot seed table order_item_statuses, with error %v\n", err)
	}
	log.Println("success seed table order_item_statuses")
}
