package usecase

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
	Validate(param []string) error
	File(filepath string)
}

type Upload interface {
	Validate(param []string) error
	File(filepathFrom string, filepathTo string)
}

type Screenshot interface {
	TakeScreenshot() error
}

type Persistence interface {
	Validate(param []string) error
	Persist(status string) error
}
