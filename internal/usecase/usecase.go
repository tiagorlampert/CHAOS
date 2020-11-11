package usecase

import "errors"

var (
	ErrRequiredParam    = errors.New("required param")
	ErrUnsupportedParam = errors.New("unsupported param")
)

type UseCase struct {
	Terminal    Terminal
	Information Information
	Download    Download
	Upload      Upload
	Screenshot  Screenshot
	Persistence Persistence
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
