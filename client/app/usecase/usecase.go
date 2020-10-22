package usecase

type UseCase struct {
	Information Information
	Terminal    Terminal
	Download    Download
	Upload      Upload
	Screenshot  Screenshot
}

type Information interface {
	Collect()
}

type Screenshot interface {
	TakeScreenshot()
}

type Download interface {
	File(data []byte)
}

type Upload interface {
	File(data []byte)
}

type Terminal interface {
	Run(cmd string)
}
