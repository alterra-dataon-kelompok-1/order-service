package model

type OrderStatus string

const (
	PendingOrder   OrderStatus = "pending"
	PaidOrder      OrderStatus = "paid"
	PreparedOrder  OrderStatus = "prepared"
	CompletedOrder OrderStatus = "completed"
	CanceledOrder  OrderStatus = "canceled"
)
