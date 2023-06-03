package forms

type Login struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Verify   string `form:"verify" binding:"required"`
}
