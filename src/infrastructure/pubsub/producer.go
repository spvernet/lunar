package pubsub

import (
	"encoding/json"

	"lunar/src/domain"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Producer struct {
	pub message.Publisher
}

func NewProducer(pub message.Publisher) *Producer {
	return &Producer{pub: pub}
}

func (p *Producer) Publish(topic string, env domain.MessageEnvelope) error {
	payload, err := json.Marshal(env)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	return p.pub.Publish(topic, msg)
}
