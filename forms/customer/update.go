package forms

type Update struct {
	Name      string  `form:"name" binding:"required"`
	ApiKey    string  `form:"apiKey" binding:"required"`
	ApiSecret string  `form:"apiSecret" binding:"required"`
	Capital   float64 `form:"capital" binding:"required"`
	Scale     float64 `form:"scale" binding:"required"`
	Level1    float64 `form:"level1" binding:"required"`
	Level2    float64 `form:"level2" binding:"required"`
	Status    string  `form:"status" binding:"required,oneof=Enable Disable"`
}
