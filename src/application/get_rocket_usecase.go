package application

import (
	"lunar/src/domain"
	"lunar/src/domain/port"
)

//go:generate moq -out get_rocket_uc_mock.go . GetRocketUCInterface
type GetRocketUCInterface interface {
	Execute(channel string) (domain.Rocket, bool, error)
}
type GetRocketUC struct {
	reader port.RocketReader
}

func NewGetRocketUC(reader port.RocketReader) GetRocketUCInterface {
	return &GetRocketUC{reader: reader}
}

func (s *GetRocketUC) Execute(channel string) (domain.Rocket, bool, error) {
	return s.reader.Get(channel)
}
