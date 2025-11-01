package application

import (
	"lunar/src/domain"

	"github.com/stretchr/testify/mock"
)

type ApplyMessageUCMock struct{ mock.Mock }

func (m *ApplyMessageUCMock) Execute(env domain.MessageEnvelope) error {
	args := m.Called(env)
	return args.Error(0)
}
