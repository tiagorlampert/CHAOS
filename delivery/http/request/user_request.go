package request

type UpdateUserPasswordRequestForm struct {
	Username    string `form:"username" json:"username,omitempty" binding:"required"`
	OldPassword string `form:"old-password" json:"old_password,omitempty" binding:"required"`
	NewPassword string `form:"new-password" json:"new_password,omitempty" binding:"required"`
}
