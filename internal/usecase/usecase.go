package usecase

type UseCase struct {
	Download   Download
	Upload     Upload
	Screenshot Screenshot
}

type Download interface {
	Validate(param []string)
	Prepare(command string)
	File(path string)
}

type Upload interface {
	Validate(param []string)
	Prepare(command string)
	File(path string)
}

type Screenshot interface {
	TakeScreenshot(input string)
}
