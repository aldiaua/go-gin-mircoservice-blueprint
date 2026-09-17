package users

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repository Repository
}

type ListResult struct {
	Users      []User
	Page       int
	Limit      int
	Total      int64
	TotalPages int
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, payload CreateUserPayload) (User, error) {
	user := User{
		ID:        uuid.NewString(),
		Name:      strings.TrimSpace(payload.Name),
		Email:     strings.ToLower(strings.TrimSpace(payload.Email)),
		CreatedAt: time.Now().UTC(),
	}
	return s.repository.Create(ctx, user)
}

func (s *Service) List(ctx context.Context, query ListUsersQuery) (ListResult, error) {
	query.Normalize()
	offset := (query.Page - 1) * query.Limit
	users, total, err := s.repository.List(ctx, offset, query.Limit)
	if err != nil {
		return ListResult{}, err
	}

	totalPages := int((total + int64(query.Limit) - 1) / int64(query.Limit))
	return ListResult{Users: users, Page: query.Page, Limit: query.Limit, Total: total, TotalPages: totalPages}, nil
}
