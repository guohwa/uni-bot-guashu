package order

type Create struct {
	CustomerID string  `form:"customer" binding:"required"`
	Symbol     string  `form:"symbol" binding:"required"`
	Action     string  `form:"action" binding:"required,oneof=OPEN CLOSE"`
	Side       string  `form:"side" binding:"required,oneof=LONG SHORT"`
	Size       float64 `form:"size" binding:"required,gt=0"`
}
