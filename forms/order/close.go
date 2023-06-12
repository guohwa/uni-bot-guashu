package order

type Close struct {
	CustomerID   string `form:"customer" binding:"required"`
	Symbol       string `form:"symbol" binding:"required"`
	PositionSide string `form:"positionSide" binding:"required,oneof=LONG SHORT"`
	PositionAmt  string `form:"positionAmt" binding:"required"`
}
