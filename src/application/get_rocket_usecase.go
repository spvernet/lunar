package application

import (
	"lunar/src/domain"
	"lunar/src/domain/port"
)

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
