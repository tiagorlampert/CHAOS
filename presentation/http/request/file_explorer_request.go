package request

type FileExplorerRequestForm struct {
	Address string `form:"address" binding:"required"`
	Path    string `form:"path"`
}
