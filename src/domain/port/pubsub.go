package port

import "lunar/src/domain"

type MessagePublisher interface {
	Publish(topic string, env domain.MessageEnvelope) error
}
