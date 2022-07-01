package repository

import (
	"context"
	"log"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, order model.Order) (*model.Order, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, order model.Order) (*model.Order, error) {
	createdOrder := order
	log.Println("createdOrder:", createdOrder)
	// createdOrder.OrderItems = nil
	err := r.db.WithContext(ctx).Create(&createdOrder).Error
	return &createdOrder, err
}
