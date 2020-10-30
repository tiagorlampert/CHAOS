package generate

import "github.com/tiagorlampert/CHAOS/internal/usecase"

type GenerateUseCase struct{}

func NewGenerateUseCase() usecase.Build {
	return &GenerateUseCase{}
}

func (g GenerateUseCase) BuildClientBinary(params []string) error {
	panic("implement me")
}
