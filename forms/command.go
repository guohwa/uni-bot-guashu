package forms

import "github.com/uncle-gua/gobinance/futures"

type Command struct {
	Action  string                   `form:"action" binding:"required"`
	Symbol  string                   `form:"symbol" binding:"required"`
	Side    futures.PositionSideType `form:"side" binding:"required,oneof=LONG SHORT"`
	Capital float64                  `form:"capital" binding:"required"`
	Size    float64                  `form:"size" binding:"required"`
	Comment string                   `form:"comment" binding:"required"`
}
