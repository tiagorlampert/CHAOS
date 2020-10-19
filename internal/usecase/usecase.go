package usecase

type UseCase struct {
	Download   Download
	Upload     Upload
	Screenshot Screenshot
}

type Download interface {
	Validate(param []string)
	Prepare(command string)
	ReceiveFile(path string)
}

type Upload interface {
	Validate(param []string)
	Prepare(command string)
	SendPath(savePath string)
	SendFile(filepath string)
}

type Screenshot interface {
	TakeScreenshot(input string)
}
