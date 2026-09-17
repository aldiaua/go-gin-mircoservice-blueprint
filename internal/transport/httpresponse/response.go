package httpresponse

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Meta struct {
	RequestID  string      `json:"request_id"`
	Timestamp  time.Time   `json:"timestamp"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type Envelope struct {
	Data  any     `json:"data,omitempty"`
	Meta  Meta    `json:"meta"`
	Error *Detail `json:"error,omitempty"`
}

type Detail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, Envelope{Data: data, Meta: meta(c)})
}

func SuccessWithPagination(c *gin.Context, status int, data any, pagination Pagination) {
	responseMeta := meta(c)
	responseMeta.Pagination = &pagination
	c.JSON(status, Envelope{Data: data, Meta: responseMeta})
}

func Failure(c *gin.Context, status int, code, message string, details any) {
	c.JSON(status, Envelope{
		Meta:  meta(c),
		Error: &Detail{Code: code, Message: message, Details: details},
	})
}

func meta(c *gin.Context) Meta {
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = uuid.NewString()
		c.Header("X-Request-ID", requestID)
	}
	return Meta{RequestID: requestID, Timestamp: time.Now().UTC()}
}

func BadRequest(c *gin.Context, code, message string, details any) {
	Failure(c, http.StatusBadRequest, code, message, details)
}

func Conflict(c *gin.Context, code, message string, details any) {
	Failure(c, http.StatusConflict, code, message, details)
}

func InternalError(c *gin.Context) {
	Failure(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an internal server error occurred", nil)
}
