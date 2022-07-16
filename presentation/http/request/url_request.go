package request

type OpenUrlRequestForm struct {
	Address string `form:"address" binding:"required"`
	URL     string `form:"url"  binding:"required"`
}
