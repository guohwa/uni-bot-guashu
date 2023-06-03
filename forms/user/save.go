package forms

type Save struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Role     string `form:"role" binding:"required,oneof=Admin User Demo"`
	Status   string `form:"status" binding:"required,oneof=Enable Disable"`
}
