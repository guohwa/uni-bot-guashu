package forms

import "github.com/uncle-gua/gobinance/futures"

type Command struct {
	Action  string                   `form:"action" binding:"required"`
	Symbol  string                   `form:"symbol" binding:"required"`
	Side    futures.PositionSideType `form:"side" binding:"required,oneof=LONG SHORT"`
	Size    float64                  `form:"size" binding:"gte=0"`
	Comment string                   `form:"comment" binding:"required"`
}
