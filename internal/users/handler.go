package users

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/example/go-gin-blueprint/internal/transport/httpresponse"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service, validate: validator.New()}
}

func (h *Handler) Create(c *gin.Context) {
	var payload CreateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		httpresponse.BadRequest(c, "INVALID_PAYLOAD", "invalid JSON payload", nil)
		return
	}
	if err := h.validate.Struct(payload); err != nil {
		httpresponse.BadRequest(c, "VALIDATION_FAILED", "validation failed", err.Error())
		return
	}

	user, err := h.service.Create(c.Request.Context(), payload)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			httpresponse.Conflict(c, "EMAIL_ALREADY_EXISTS", "email already exists", nil)
			return
		}
		httpresponse.InternalError(c)
		return
	}
	httpresponse.Success(c, 201, NewUserResponse(user))
}

func (h *Handler) List(c *gin.Context) {
	var query ListUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httpresponse.BadRequest(c, "INVALID_QUERY", "invalid pagination parameters", nil)
		return
	}
	if err := h.validate.Struct(query); err != nil {
		httpresponse.BadRequest(c, "VALIDATION_FAILED", "invalid pagination parameters", err.Error())
		return
	}

	result, err := h.service.List(c.Request.Context(), query)
	if err != nil {
		httpresponse.InternalError(c)
		return
	}
	httpresponse.SuccessWithPagination(c, 200, NewUserResponses(result.Users), httpresponse.Pagination{
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: result.TotalPages,
	})
}
