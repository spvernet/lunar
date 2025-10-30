package port

import "lunar/src/domain"

type MessageWriter interface {
	Apply(env domain.MessageEnvelope) error
}

type RocketReader interface {
	Get(channel string) (domain.Rocket, bool, error)
	List(sortBy string, order string) ([]domain.Rocket, error)
}

type Persistence interface {
	MessageWriter
	RocketReader
}
