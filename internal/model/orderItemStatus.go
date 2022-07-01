package model

type OrderItemStatus string

const (
	Pending   OrderItemStatus = "pending"
	Paid      OrderItemStatus = "paid"
	Prepared  OrderItemStatus = "prepared"
	Delivered OrderItemStatus = "delivered"
)
