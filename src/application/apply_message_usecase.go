package application

import (
	"lunar/src/domain"
	"lunar/src/domain/port"
)

//go:generate moq -out apply_message_uc_mock.go . ApplyMessageUCInterface
type ApplyMessageUCInterface interface {
	Execute(env domain.MessageEnvelope) error
}
type ApplyMessageUC struct {
	writer port.MessageWriter
}

func NewApplyMessageUC(writer port.MessageWriter) ApplyMessageUCInterface {
	return &ApplyMessageUC{writer: writer}
}

func (uc *ApplyMessageUC) Execute(env domain.MessageEnvelope) error {
	return uc.writer.Apply(env)
}
