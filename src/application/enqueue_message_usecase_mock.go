package application

import (
	"lunar/src/domain"

	"github.com/stretchr/testify/mock"
)

type EnqueueMessageUCMock struct{ mock.Mock }

func (m *EnqueueMessageUCMock) Execute(env domain.MessageEnvelope) error {
	args := m.Called(env)
	return args.Error(0)
}
