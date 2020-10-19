package usecase

type UseCase struct {
	Download    Download
	Upload      Upload
	Screenshot  Screenshot
	Information Information
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

type Information interface {
	Collect()
}
