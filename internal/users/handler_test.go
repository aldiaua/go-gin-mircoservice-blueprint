package users

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestRouter(repository Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router.Group("/api/v1"), NewHandler(NewService(repository)))
	return router
}

func TestHandlerCreateRejectsInvalidPayload(t *testing.T) {
	router := newTestRouter(&repositoryStub{})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(`{"name":"A","email":"bad"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Contains(t, response.Body.String(), "validation failed")
	require.NotEmpty(t, response.Header().Get("X-Request-ID"))
	require.Contains(t, response.Body.String(), "request_id")
}

func TestHandlerCreateReturnsCreatedUser(t *testing.T) {
	repository := &repositoryStub{}
	router := newTestRouter(repository)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(`{"name":"Ada Lovelace","email":"ada@example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	require.Contains(t, response.Body.String(), "ada@example.com")
	require.NotEmpty(t, response.Header().Get("X-Request-ID"))
	require.Contains(t, response.Body.String(), "\"meta\"")
	require.NotEmpty(t, repository.created.ID)
}

func TestHandlerCreateReturnsConflictForDuplicateEmail(t *testing.T) {
	router := newTestRouter(&repositoryStub{err: ErrEmailAlreadyExists})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(`{"name":"Ada Lovelace","email":"ada@example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusConflict, response.Code)
	require.Contains(t, response.Body.String(), "EMAIL_ALREADY_EXISTS")
}

func TestHandlerListAppliesPagination(t *testing.T) {
	router := newTestRouter(&repositoryStub{users: []User{{ID: "user-1", Name: "Ada", Email: "ada@example.com"}}})
	request := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=2&limit=10", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Body.String(), `"page":2`)
	require.Contains(t, response.Body.String(), `"limit":10`)
}
