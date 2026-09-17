package httpmiddleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRequestIDPreservesValidUUID(t *testing.T) {
	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) { c.Status(204) })
	request := httptest.NewRequest("GET", "/", nil)
	expected := uuid.NewString()
	request.Header.Set(requestIDHeader, expected)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, expected, response.Header().Get(requestIDHeader))
}

func TestRequestIDReplacesInvalidValue(t *testing.T) {
	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) { c.Status(204) })
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set(requestIDHeader, "invalid-request-id")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	requestID := response.Header().Get(requestIDHeader)
	_, err := uuid.Parse(requestID)
	require.NoError(t, err)
	require.NotEqual(t, "invalid-request-id", requestID)
}
