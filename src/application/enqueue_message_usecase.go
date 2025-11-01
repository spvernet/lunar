package application

import (
	"lunar/src/domain"
	"lunar/src/domain/port"
)

type EnqueueMessageUCInterface interface {
	Execute(env domain.MessageEnvelope) error
}
type EnqueueMessageUC struct {
	pub   port.MessagePublisher
	topic string
}

func NewEnqueueMessageUC(pub port.MessagePublisher, topic string) EnqueueMessageUCInterface {
	return &EnqueueMessageUC{pub: pub, topic: topic}
}

func (uc *EnqueueMessageUC) Execute(env domain.MessageEnvelope) error {
	return uc.pub.Publish(uc.topic, env)
}
