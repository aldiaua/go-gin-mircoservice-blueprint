package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/example/go-gin-blueprint/internal/config"
)

func testServer() *HTTPServer {
	return NewHTTPServer(config.Config{Port: "0"}, &gorm.DB{}, logrus.New())
}

func TestHealthRouteReturnsStandardResponse(t *testing.T) {
	app := testServer()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	app.httpServer.Handler.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, response.Header().Get("X-Request-ID"))
	require.Contains(t, response.Body.String(), `"status":"ok"`)
	require.Contains(t, response.Body.String(), `"request_id"`)
}

func TestRunStopsWhenContextIsCancelled(t *testing.T) {
	app := testServer()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.NoError(t, app.Run(ctx))
}
