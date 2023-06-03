package forms

type Update struct {
	Password string `form:"password" binding:""`
	Role     string `form:"role" binding:"required,oneof=Admin User Demo"`
	Status   string `form:"status" binding:"required,oneof=Enable Disable"`
}
