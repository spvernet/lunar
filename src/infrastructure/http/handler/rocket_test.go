package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"lunar/src/application"
	"lunar/src/domain"
	h "lunar/src/infrastructure/http/handler"
)

const (
	pathGetOne  = "/api/rockets/:channel"
	pathList    = "/api/rockets"
	urlGetOne   = "/api/rockets/%s"
	urlList     = "/api/rockets?sort=%s&order=%s"
	contentType = "application/json"
)

// -------- helpers --------

func newRocketsRouter(t *testing.T, getUC application.GetRocketUCInterface, listUC application.ListRocketsUCInterface) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	hdl := h.NewRockets(getUC, listUC)
	r.GET(pathGetOne, hdl.GetOne)
	r.GET(pathList, hdl.List)
	return r
}

func doGET(r *gin.Engine, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set(headerCT, contentType)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// -------- tests: GetOne --------

func TestRockets_GetOne_HappyPath_Returns200(t *testing.T) {
	getMock := &application.GetRocketUCMock{}
	listMock := &application.ListRocketsUCMock{} // not used here

	ch := "abc"
	rc := domain.Rocket{Channel: ch}

	getMock.
		On("Execute", ch).
		Return(rc, true, nil).
		Once()

	r := newRocketsRouter(t, getMock, listMock)
	w := doGET(r, fmt.Sprintf(urlGetOne, ch))

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	getMock.AssertExpectations(t)
	getMock.AssertNumberOfCalls(t, "Execute", 1)
}

func TestRockets_GetOne_NotFound_Returns404(t *testing.T) {
	getMock := &application.GetRocketUCMock{}
	listMock := &application.ListRocketsUCMock{}

	ch := "missing"

	getMock.
		On("Execute", ch).
		Return(domain.Rocket{}, false, nil).
		Once()

	r := newRocketsRouter(t, getMock, listMock)
	w := doGET(r, fmt.Sprintf(urlGetOne, ch))

	require.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
	getMock.AssertExpectations(t)
	getMock.AssertNotCalled(t, "Execute", "other") // sanity check
}

func TestRockets_GetOne_InternalError_Returns500(t *testing.T) {
	getMock := &application.GetRocketUCMock{}
	listMock := &application.ListRocketsUCMock{}

	ch := "abc"

	getMock.
		On("Execute", ch).
		Return(domain.Rocket{}, false, fmt.Errorf("db down")).
		Once()

	r := newRocketsRouter(t, getMock, listMock)
	w := doGET(r, fmt.Sprintf(urlGetOne, ch))

	require.Equal(t, http.StatusInternalServerError, w.Code, w.Body.String())
	getMock.AssertExpectations(t)
}

// -------- tests: List --------

func TestRockets_List_HappyPath_Returns200(t *testing.T) {
	getMock := &application.GetRocketUCMock{}
	listMock := &application.ListRocketsUCMock{}

	sortBy := h.SortBySpeed
	order := h.OrderDesc
	items := []domain.Rocket{{Channel: "a"}, {Channel: "b"}}

	listMock.
		On("Execute", sortBy, order).
		Return(items, nil).
		Once()

	r := newRocketsRouter(t, getMock, listMock)
	w := doGET(r, fmt.Sprintf(urlList, sortBy, order))

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	listMock.AssertExpectations(t)
	listMock.AssertNumberOfCalls(t, "Execute", 1)
}

func TestRockets_List_InvalidSort_Returns400_AndDoesNotCallUC(t *testing.T) {
	getMock := &application.GetRocketUCMock{}
	listMock := &application.ListRocketsUCMock{}

	r := newRocketsRouter(t, getMock, listMock)
	w := doGET(r, fmt.Sprintf(urlList, "oops", h.OrderAsc))

	require.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	listMock.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything)
}

func TestRockets_List_InvalidOrder_Returns400_AndDoesNotCallUC(t *testing.T) {
	getMock := &application.GetRocketUCMock{}
	listMock := &application.ListRocketsUCMock{}

	r := newRocketsRouter(t, getMock, listMock)
	w := doGET(r, fmt.Sprintf(urlList, h.SortByChannel, "down"))

	require.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	listMock.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything)
}

func TestRockets_List_InternalError_Returns500(t *testing.T) {
	getMock := &application.GetRocketUCMock{}
	listMock := &application.ListRocketsUCMock{}

	sortBy := h.SortByChannel
	order := h.OrderAsc

	listMock.
		On("Execute", sortBy, order).
		Return(nil, fmt.Errorf("repo error")).
		Once()

	r := newRocketsRouter(t, getMock, listMock)
	w := doGET(r, fmt.Sprintf(urlList, sortBy, order))

	require.Equal(t, http.StatusInternalServerError, w.Code, w.Body.String())
	listMock.AssertExpectations(t)
}
