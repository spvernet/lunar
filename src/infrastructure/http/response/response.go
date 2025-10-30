package response

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

func WriteErrorResponse(c *gin.Context, code int, err error) {
	c.Writer.WriteHeader(code)
	c.Writer.Header().Set(ContentType, ApplicationJson)
	_, _ = c.Writer.WriteString(
		fmt.Sprintf("{\"code\":%d,\"message\":\"%s\"}", code, err.Error()),
	)
}

func WriteEmptyResponse(c *gin.Context, code int) {
	c.Writer.WriteHeader(code)
	c.Writer.Header().Set(ContentType, ApplicationJson)
}

func WriteJSONResponse(c *gin.Context, code int, msg any) {
	c.Writer.WriteHeader(code)
	c.Writer.Header().Set(ContentType, ApplicationJson)
	c.JSON(code, msg)
}
