package repository

import (
	"context"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetOrders(ctx context.Context, payload *dto.GetRequest) (*[]model.Order, *dto.PaginationInfo, error)
	Create(ctx context.Context, order model.Order) (*model.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetOrders(ctx context.Context, payload *dto.GetRequest) (*[]model.Order, *dto.PaginationInfo, error) {
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

func (r *repository) GetOrderByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	var data model.Order
	err := r.db.WithContext(ctx).First(&data, id).Error

	return &data, err
}
