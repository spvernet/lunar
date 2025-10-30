package handler

import (
	httpresponse "lunar/src/infrastructure/http/response"
	"net/http"

	"lunar/src/application"
	"lunar/src/domain"
	"lunar/src/domain/validator"

	"github.com/gin-gonic/gin"
)

type Messages struct {
	enqueue   application.EnqueueMessageUCInterface
	validator validator.Validator
}

func NewMessages(
	enqueue application.EnqueueMessageUCInterface,
	v validator.Validator,
) *Messages {
	return &Messages{enqueue: enqueue, validator: v}
}

func (h *Messages) Handle(c *gin.Context) {
	var env domain.MessageEnvelope
	if err := c.ShouldBindJSON(&env); err != nil {
		httpresponse.WriteErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := h.validator.ValidateEnvelope(env); err != nil {
		httpresponse.WriteErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := h.validator.ValidatePayload(env.Metadata.MessageType, env.Message); err != nil {
		httpresponse.WriteErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := h.enqueue.Execute(env); err != nil {
		httpresponse.WriteErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	httpresponse.WriteEmptyResponse(c, http.StatusAccepted)
}
