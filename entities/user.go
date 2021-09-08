package entities

type User struct {
	DBModel
	Username string `form:"username" json:"username,omitempty" binding:"required"`
	Password string `form:"password" json:"password,omitempty" binding:"required"`
}
