package usecase

import "errors"

var (
	ErrRequiredParam    = errors.New("required param")
	ErrUnsupportedParam = errors.New("unsupported param")
	ErrSendingRequest   = errors.New("error sending request")
	ErrReadingResponse  = errors.New("error reading response")
)

type UseCase struct {
	Terminal    Terminal
	Information Information
	Download    Download
	Upload      Upload
	Screenshot  Screenshot
	Persistence Persistence
	OpenURL     OpenURL
}

type Information interface {
	Collect()
}

type Terminal interface {
	Run(cmd string)
}

type Download interface {
	Validate(params []string) error
	File(filepath string)
}

type Upload interface {
	Validate(params []string) error
	File(filepathFrom string, filepathTo string)
}

type Screenshot interface {
	TakeScreenshot() error
}

type Persistence interface {
	Validate(params []string) error
	Persist(status string) error
}

type Build interface {
	BuildClientBinary(params []string) error
}

type OpenURL interface {
	Open(url []string) error
}
