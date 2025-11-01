package pubsub

import (
	"context"
	"encoding/json"

	"lunar/src/application"
	"lunar/src/domain"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

type Consumer struct {
	sub     message.Subscriber
	log     *zap.Logger
	applyUC application.ApplyMessageUCInterface
	topic   string
}

func NewConsumer(
	sub message.Subscriber,
	log *zap.Logger,
	applyUC application.ApplyMessageUCInterface,
	topic string,
) *Consumer {
	return &Consumer{
		sub:     sub,
		log:     log,
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
			if err = json.Unmarshal(msg.Payload, &env); err != nil {
				// NACK: Watermill gochannel reentrega (o se pierde según config);
				// aquí ACK para no bloquear; loguea el error real.
				c.log.Error(err.Error())
				msg.Ack()
				continue
			}
			if err = c.applyUC.Execute(env); err != nil {
				// Puedes Nack() si quieres reintentar; en gochannel no hay persistencia.
				c.log.Error(err.Error())
				msg.Ack()
				continue
			}
			msg.Ack()
		}
	}()
	return nil
}
