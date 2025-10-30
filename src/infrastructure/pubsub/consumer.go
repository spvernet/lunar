package pubsub

import (
	"context"
	"encoding/json"

	"lunar/src/application"
	"lunar/src/domain"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

// Consumer: suscribe a topic, deserializa y delega en ApplyMessageUC.
type Consumer struct {
	sub         message.Subscriber
	logger      watermill.LoggerAdapter
	applyUC     application.ApplyMessageUCInterface
	topic       string
	workerGroup string // optional: para logs/metrics
}

func NewConsumer(
	sub message.Subscriber,
	logger watermill.LoggerAdapter,
	applyUC application.ApplyMessageUCInterface,
	topic string,
) *Consumer {
	return &Consumer{
		sub:     sub,
		logger:  logger,
		applyUC: applyUC,
		topic:   topic,
	}
}

func (c *Consumer) Subscribe(ctx context.Context) error {
	msgs, err := c.sub.Subscribe(ctx, c.topic)
	if err != nil {
		return err
	}
	go func() {
		for msg := range msgs {
			var env domain.MessageEnvelope
			if err := json.Unmarshal(msg.Payload, &env); err != nil {
				// NACK: Watermill gochannel reentrega (o se pierde según config);
				// aquí ACK para no bloquear; loguea el error real.
				c.logger.Error("failed to unmarshal", err, watermill.LogFields{"topic": c.topic})
				msg.Ack()
				continue
			}
			if err := c.applyUC.Execute(env); err != nil {
				// Puedes Nack() si quieres reintentar; en gochannel no hay persistencia.
				c.logger.Error("apply failed", err, watermill.LogFields{"topic": c.topic})
				msg.Ack()
				continue
			}
			msg.Ack()
		}
	}()
	return nil
}
