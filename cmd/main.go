package main

import (
	"context"
	"lunar/src/domain/validator"

	"lunar/src/application"
	"lunar/src/infrastructure/http/handler"
	"lunar/src/infrastructure/http/routes"
	"lunar/src/infrastructure/persistence"
	"lunar/src/infrastructure/pubsub"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const topicMessages = "rockets.messages"

func mustSucceed[C any](h C, err error) C {
	if err != nil {
		panic(err)
	}
	return h
}

func main() {
	mem := persistence.NewMemoryStore()
	logger := mustSucceed(zap.NewDevelopment())
	ctx := context.Background()

	channel := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	// Producer & Consumer
	producer := pubsub.NewProducer(channel)
	applyUC := application.NewApplyMessageUC(mem)
	consumer := pubsub.NewConsumer(channel, logger, applyUC, topicMessages)

	// Arranca el consumer
	if err := consumer.Subscribe(ctx); err != nil {
		logger.Fatal("failed to subscribe", zap.Error(err))
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

	logger.Info("listening on :8088")
	if err := r.Run(":8088"); err != nil {
		logger.Fatal("failed to run", zap.Error(err))
	}
}
