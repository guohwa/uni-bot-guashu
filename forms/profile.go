package forms

type Profile struct {
	ApiKey    string `form:"apiKey" binding:"required"`
	ApiSecret string `form:"apiSecret" binding:"required"`
}
