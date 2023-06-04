package forms

type Save struct {
	Name      string  `form:"name" binding:"required"`
	ApiKey    string  `form:"apiKey" binding:"required"`
	ApiSecret string  `form:"apiSecret" binding:"required"`
	Base      float64 `form:"base" binding:"required"`
	Capital   float64 `form:"capital" binding:"required"`
	Ratio     float64 `form:"ratio" binding:"required"`
	Level1    float64 `form:"level1" binding:"required"`
	Level2    float64 `form:"level2" binding:"required"`
	Mode      string  `form:"mode" binding:"required,oneof=Capital Equity"`
	Status    string  `form:"status" binding:"required,oneof=Enable Disable"`
}
