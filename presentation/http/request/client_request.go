package request

type GenerateClientRequestForm struct {
	Address   string `form:"address"`
	Port      string `form:"port"`
	OSTarget  string `form:"os_target"`
	Filename  string `form:"filename"`
	RunHidden string `form:"run_hidden"`
}
