package repository

import (
	"context"
	"log"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetOrder(ctx context.Context, payload *dto.GetRequest) (*[]model.Order, *dto.PaginationInfo, error)
	Create(ctx context.Context, order model.Order) (*model.Order, error)
	UpdateOrderByID(ctx context.Context, id uuid.UUID, data map[string]interface{}) error
	UpdateOrderByIDWithModel(ctx context.Context, id uuid.UUID, data *model.Order) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetOrder(ctx context.Context, payload *dto.GetRequest) (*[]model.Order, *dto.PaginationInfo, error) {
	var orders []model.Order
	var count int64

	query := r.db.WithContext(ctx).Model(&model.Order{}).Preload("OrderItems")
	countQuery := query
	if countQuery.Error != nil {
		return nil, nil, query.Error
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, nil, err
	}

	limit, offset := payload.Pagination.GetLimitOffset()
	err := query.Limit(limit).Offset(offset).Find(&orders).Error

	return &orders, payload.Pagination.CheckInfoPagination(count), err
}

func (r *repository) Create(ctx context.Context, order model.Order) (*model.Order, error) {
	err := r.db.WithContext(ctx).Create(&order).Error
	return &order, err
}

func (r *repository) UpdateOrderByID(ctx context.Context, id uuid.UUID, data map[string]interface{}) error {
	log.Println("data in repo:", data)
	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&model.Order{}).Updates(data).Error
	return err
}

func (r *repository) UpdateOrderByIDWithModel(ctx context.Context, id uuid.UUID, data *model.Order) error {
	log.Println("data in repo:", data)

	if data.OrderItems != nil {
		qty := data.OrderItems[0].Quantity
		log.Println("qty:", qty)
	}

	// err := r.db.WithContext(ctx).Where("id = ?", id).Model(&model.Order{}).Updates(data).Error
	err := r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Model(&model.Order{}).Updates(&data).Error
	return err
}

func (r *repository) UpdateOrderItemByID(ctx context.Context, id uuid.UUID, data map[string]interface{}) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Model(&model.OrderItem{}).Preload("OrderItems").Updates(data).Error
	return err
}
