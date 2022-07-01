package database

import (
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		model.Order{},
		model.OrderItem{},
		model.OrderItemStatus{},
		model.OrderStatus{},
	)
}

func DropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		model.OrderItem{},
		model.OrderItemStatus{},
		model.OrderStatus{},
		model.Order{},
	)
}
