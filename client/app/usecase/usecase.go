package usecase

import "errors"

var ErrUnsupportedPlatform = errors.New("unsupported platform")

type UseCase struct {
	Information Information
	Terminal    Terminal
	Download    Download
	Upload      Upload
	Screenshot  Screenshot
	Persistence Persistence
	OpenURL     OpenURL
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

type Persistence interface {
	Persist(data []byte)
}

type OpenURL interface {
	Open(url string)
}
