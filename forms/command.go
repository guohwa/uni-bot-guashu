package forms

type Command struct {
	Action  string  `form:"action" binding:"required"`
	Symbol  string  `form:"symbol" binding:"required"`
	Side    string  `form:"side" binding:"required"`
	Capital float64 `form:"capital" binding:"required"`
	Size    float64 `form:"size" binding:"required"`
	Comment string  `form:"comment" binding:"required"`
}
