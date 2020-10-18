package usecase

type UseCase struct {
	Download   Download
	Upload     Upload
	Screenshot Screenshot
}

type Download interface {
	File()
}

type Upload interface {
	File()
}

type Screenshot interface {
	TakeScreenshot()
}
