package handler

import (
	"lunar/src/application"
	"lunar/src/infrastructure/http/httperror"
	"lunar/src/infrastructure/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	KeySort         = "sort"
	KeyOrder        = "order"
	SortByChannel   = "channel"
	SortBySpeed     = "speed"
	SortByUpdatedAt = "updated_at"

	OrderAsc  = "asc"
	OrderDesc = "desc"
)

type Rockets struct {
	get  application.GetRocketUCInterface
	list application.ListRocketsUCInterface
}

func NewRockets(get application.GetRocketUCInterface, list application.ListRocketsUCInterface) *Rockets {
	return &Rockets{get: get, list: list}
}

func (h *Rockets) GetOne(c *gin.Context) {
	ch := c.Param("channel")
	rocket, ok, err := h.get.Execute(ch)
	if err != nil {
		response.WriteErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		response.WriteErrorResponse(c, http.StatusNotFound, httperror.ErrChannelNotFound)
		return
	}
	response.WriteJSONResponse(c, http.StatusOK, rocket)
}

func (h *Rockets) List(c *gin.Context) {

	sortBy := c.DefaultQuery(KeySort, SortByChannel)
	order := c.DefaultQuery(KeyOrder, OrderAsc)
	if sortBy != SortByChannel && sortBy != SortBySpeed && sortBy != SortByUpdatedAt {
		response.WriteErrorResponse(
			c,
			http.StatusBadRequest,
			httperror.ErrInvalidSort,
		)
		return
	}
	if order != OrderAsc && order != OrderDesc {
		response.WriteErrorResponse(c, http.StatusBadRequest, httperror.ErrInvalidOrder)
		return
	}
	items, err := h.list.Execute(sortBy, order)
	if err != nil {
		response.WriteErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	response.WriteJSONResponse(c, http.StatusOK, items)
}
