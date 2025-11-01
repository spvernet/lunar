package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"lunar/src/application"
	"lunar/src/domain/validator"
	h "lunar/src/infrastructure/http/handler"
)

const (
	pathMessages = "/messages"
	headerCT     = "Content-Type"
	ctJSON       = "application/json"
)

func TestMessages_InvalidJSON_Returns400_AndDoesNotCallEnqueue(t *testing.T) {
	ucMock := &application.EnqueueMessageUCMock{}
	r := newRouter(t, ucMock)

	w := postJSON(r, pathMessages, `{"metadata":{`) // JSON truncado

	require.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	ucMock.AssertNotCalled(t, "Execute", mock.Anything)
}

func TestMessages_InvalidEnvelope_Returns400_AndDoesNotCallEnqueue(t *testing.T) {
	ucMock := &application.EnqueueMessageUCMock{}
	r := newRouter(t, ucMock)

	now := time.Now().Format(time.RFC3339Nano)
	body := fmt.Sprintf(`{
	  "metadata":{"channel":"","messageNumber":1,"messageTime":"%s","messageType":"RocketLaunched"},
	  "message":{"type":"Falcon-9","launchSpeed":100,"mission":"M1"}
	}`, now)

	w := postJSON(r, pathMessages, body)

	require.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	ucMock.AssertNotCalled(t, "Execute", mock.Anything)
}

func TestMessages_InvalidPayload_Returns400_AndDoesNotCallEnqueue(t *testing.T) {
	ucMock := &application.EnqueueMessageUCMock{}
	r := newRouter(t, ucMock)

	now := time.Now().Format(time.RFC3339Nano)
	body := fmt.Sprintf(`{
	  "metadata":{"channel":"c1","messageNumber":2,"messageTime":"%s","messageType":"RocketSpeedIncreased"},
	  "message":{"by":0}
	}`, now) // by <= 0 -> inválido

	w := postJSON(r, pathMessages, body)

	require.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	ucMock.AssertNotCalled(t, "Execute", mock.Anything)
}

func TestMessages_EnqueueError_Returns500_AndCallsEnqueueOnce(t *testing.T) {
	ucMock := &application.EnqueueMessageUCMock{}
	ucMock.
		On("Execute", mock.AnythingOfType("domain.MessageEnvelope")).
		Return(fmt.Errorf("boom")).
		Once()

	r := newRouter(t, ucMock)

	now := time.Now().Format(time.RFC3339Nano)
	body := fmt.Sprintf(`{
	  "metadata":{"channel":"c1","messageNumber":1,"messageTime":"%s","messageType":"RocketLaunched"},
	  "message":{"type":"Falcon-9","launchSpeed":100,"mission":"M1"}
	}`, now)

	w := postJSON(r, pathMessages, body)

	require.Equal(t, http.StatusInternalServerError, w.Code, w.Body.String())
	ucMock.AssertNumberOfCalls(t, "Execute", 1)
	ucMock.AssertExpectations(t)
}

func TestMessages_HappyPath_Returns202_AndCallsEnqueueOnce(t *testing.T) {
	ucMock := &application.EnqueueMessageUCMock{}
	ucMock.
		On("Execute", mock.AnythingOfType("domain.MessageEnvelope")).
		Return(nil).
		Once()

	r := newRouter(t, ucMock)

	now := time.Now().Format(time.RFC3339Nano)
	body := fmt.Sprintf(`{
	  "metadata":{"channel":"c-ok","messageNumber":1,"messageTime":"%s","messageType":"RocketLaunched"},
	  "message":{"type":"Falcon-9","launchSpeed":100,"mission":"M1"}
	}`, now)

	w := postJSON(r, pathMessages, body)

	require.Equal(t, http.StatusAccepted, w.Code, w.Body.String())
	ucMock.AssertNumberOfCalls(t, "Execute", 1)
	ucMock.AssertExpectations(t)
}

func newRouter(t *testing.T, uc application.EnqueueMessageUCInterface) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	v := validator.New()
	msgHandler := h.NewMessages(uc, v)

	r.POST(pathMessages, msgHandler.Handle)
	return r
}

func postJSON(r *gin.Engine, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set(headerCT, ctJSON)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
