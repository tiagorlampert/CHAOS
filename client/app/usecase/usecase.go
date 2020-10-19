package usecase

type UseCase struct {
	Information Information
	Screenshot  Screenshot
	Download    Download
	Upload      Upload
	Terminal    Terminal
}

type Information interface {
	Collect()
}

type Screenshot interface {
	TakeScreenshot()
}

type Download interface {
	File()
}

type Upload interface {
	ValidatePath()
	StoreFile()
}

type Terminal interface {
	Run(cmd string)
}
