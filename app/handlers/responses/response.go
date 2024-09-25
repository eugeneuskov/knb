package responses

import (
	"errors"
	"github.com/gin-gonic/gin"
	customErrors "knb/app/errors"
	"net/http"
	"strings"
)

type Response struct {
}

func NewResponse() *Response {
	return &Response{}
}

type errorResponse struct {
	Message string `json:"message"`
}

func (r *Response) NewErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{
		Message: prepareErrorMessage(message),
	})
}

func (r *Response) ParseError(c *gin.Context, err error) {
	var statusCode int
	var message string

	var repositoryUniqueViolationError *customErrors.UniqueViolationError
	var wrongLoginError *customErrors.WrongLoginError

	if errors.As(err, &repositoryUniqueViolationError) {
		statusCode = http.StatusConflict
		message = repositoryUniqueViolationError.Error()
	} else if errors.As(err, &wrongLoginError) {
		statusCode = http.StatusUnauthorized
		message = wrongLoginError.Error()
	} else {
		statusCode = http.StatusInternalServerError
		message = err.Error()
	}

	r.NewErrorResponse(c, statusCode, message)
}

func (r *Response) NewOkResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

func prepareErrorMessage(message string) string {
	messageParts := strings.Split(message, ":")

	return messageParts[len(messageParts)-1]
}
