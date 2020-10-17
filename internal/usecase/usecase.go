package usecase

type UseCase struct {
	Screenshot Screenshot
	Download   Download
	Upload     Upload
}

type Screenshot interface {
	TakeScreenshot(input string)
}

type Download interface {
}

type Upload interface {
}
