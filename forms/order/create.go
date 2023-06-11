package order

type Create struct {
	CustomerID string  `form:"customer" binding:"required"`
	Symbol     string  `form:"symbol" binding:"required"`
	Action     string  `form:"action" binding:"required"`
	Side       string  `form:"side" binding:"required"`
	Size       float64 `form:"size" binding:"required"`
}
