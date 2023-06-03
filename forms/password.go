package forms

type Password struct {
	Password        string `form:"password" binding:"required"`
	NewPassword     string `form:"password1" binding:"required,eqfield=ConfirmPassword"`
	ConfirmPassword string `form:"password2" binding:"required,eqfield=NewPassword"`
}
