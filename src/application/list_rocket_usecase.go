package application

import (
	"lunar/src/domain"
	"lunar/src/domain/port"
)

type ListRocketsUCInterface interface {
	Execute(sortBy, order string) ([]domain.Rocket, error)
}

type ListRocketsUC struct {
	reader port.RocketReader
}

func NewListRocketsUC(reader port.RocketReader) ListRocketsUCInterface {
	return &ListRocketsUC{reader: reader}
}

func (s *ListRocketsUC) Execute(sortBy, order string) ([]domain.Rocket, error) {
	return s.reader.List(sortBy, order)
}
