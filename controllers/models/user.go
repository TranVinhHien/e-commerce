package controllers_model

type LoginParams struct {
	Username string `form:"username" json:"username" binding:"required,alphanum"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
}
type LogOutParams struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}
type RegisterParams struct {
	Username string `form:"username" json:"username" binding:"required,alphanum"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
	FullName string `form:"full_name" json:"full_name" binding:"required"`
}
