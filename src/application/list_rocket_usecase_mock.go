package application

import (
	"lunar/src/domain"

	"github.com/stretchr/testify/mock"
)

type ListRocketsUCMock struct{ mock.Mock }

func (m *ListRocketsUCMock) Execute(sortBy, order string) ([]domain.Rocket, error) {
	args := m.Called(sortBy, order)

	var items []domain.Rocket
	if v, ok := args.Get(0).([]domain.Rocket); ok {
		items = v
	}
	return items, args.Error(1)
}
