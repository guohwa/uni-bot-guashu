package forms

type Save struct {
	Name      string  `form:"name" binding:"required"`
	ApiKey    string  `form:"apiKey" binding:"required"`
	ApiSecret string  `form:"apiSecret" binding:"required"`
	Capital   float64 `form:"capital" binding:"gte=0"`
	Scale     float64 `form:"scale" binding:"gt=0"`
	Level1    float64 `form:"level1" binding:"gt=0"`
	Level2    float64 `form:"level2" binding:"gt=0"`
	Status    string  `form:"status" binding:"required,oneof=Enable Disable"`
}
