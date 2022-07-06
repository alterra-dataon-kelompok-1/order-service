package order

import (
	"testing"

	"github.com/alterra-dataon-kelompok-1/order-service/internal/dto"
	"github.com/alterra-dataon-kelompok-1/order-service/internal/model"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, float32(0), res)

	// Test for 1 item
	orderItem = []model.OrderItem{
		{Price: 1_000.05, Quantity: 3},
	}
	res = sumItemPrice(orderItem)
	assert.Equal(t, float32(3_000.15), res)

	// Test for more than 1 item
	orderItem = []model.OrderItem{
		{Price: 1_000.0, Quantity: 1},
		{Price: 1_000.02, Quantity: 1},
	}
	res = sumItemPrice(orderItem)
	assert.Equal(t, float32(2_000.02), res)
}
