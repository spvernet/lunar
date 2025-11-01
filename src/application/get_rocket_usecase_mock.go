package application

import (
	"lunar/src/domain"

	"github.com/stretchr/testify/mock"
)

type GetRocketUCMock struct{ mock.Mock }

func (m *GetRocketUCMock) Execute(channel string) (domain.Rocket, bool, error) {
	args := m.Called(channel)

	var r domain.Rocket
	if v, ok := args.Get(0).(domain.Rocket); ok {
		r = v
	}
	return r, args.Bool(1), args.Error(2)
}
