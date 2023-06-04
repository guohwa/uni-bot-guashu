package forms

type Order struct {
	CustomerID string `form:"customer" binding:"required"`
	Symbol     string `form:"symbol" binding:"required"`
	OrderID    int64  `form:"orderId" binding:"required"`
}
