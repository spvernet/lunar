package main

import (
	"context"
	"log"
	"lunar/src/domain/validator"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gin-gonic/gin"

	"lunar/src/application"
	"lunar/src/infrastructure/http/handler"
	"lunar/src/infrastructure/http/routes"
	"lunar/src/infrastructure/persistence"
	"lunar/src/infrastructure/pubsub"
)

const topicMessages = "rockets.messages"

func main() {
	// Store in-memory
	mem := persistence.NewMemoryStore()

	// Logger Watermill
	wlog := watermill.NewStdLogger(false, false)

	// Canal in-memory (puedes cambiar a Kafka, Rabbit, etc.)
	channel := gochannel.NewGoChannel(gochannel.Config{}, wlog)

	// Producer & Consumer
	producer := pubsub.NewProducer(channel, wlog)
	applyUC := application.NewApplyMessageUC(mem)
	consumer := pubsub.NewConsumer(channel, wlog, applyUC, topicMessages)

	// Arranca el consumer
	if err := consumer.Subscribe(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Usecases para HTTP
	enqueueUC := application.NewEnqueueMessageUC(producer, topicMessages)
	getUC := application.NewGetRocketUC(mem)
	listUC := application.NewListRocketsUC(mem)

	// Handlers HTTP
	v := validator.New()
	msgHandler := handler.NewMessages(enqueueUC, v)
	rockHandler := handler.NewRockets(getUC, listUC)

	// Router Gin
	r := gin.Default()
	r.POST(routes.PostMessagesPath, msgHandler.Handle)

	protected := r.Group(routes.ApiGroup)
	{
		protected.GET(routes.ListRocketsPath, rockHandler.List)
		protected.GET(routes.GetRocketPath, rockHandler.GetOne)
	}

	log.Println("listening on :8088")
	if err := r.Run(":8088"); err != nil {
		log.Fatal(err)
	}
}
