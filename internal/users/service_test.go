package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type repositoryStub struct {
	created User
	users   []User
	err     error
}

func (s *repositoryStub) Create(_ context.Context, user User) (User, error) {
	if s.err != nil {
		return User{}, s.err
	}
	s.created = user
	return user, nil
}

func (s *repositoryStub) List(_ context.Context, _, _ int) ([]User, int64, error) {
	return s.users, int64(len(s.users)), nil
}

func TestServiceCreateNormalizesUserInput(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository)

	user, err := service.Create(context.Background(), CreateUserPayload{
		Name:  "  Ada Lovelace  ",
		Email: " ADA@EXAMPLE.COM ",
	})

	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	require.Equal(t, "Ada Lovelace", repository.created.Name)
	require.Equal(t, "ada@example.com", repository.created.Email)
	require.False(t, user.CreatedAt.IsZero())
}

func TestServiceListDelegatesToRepository(t *testing.T) {
	expected := []User{{ID: "user-1", Name: "Ada", Email: "ada@example.com"}}
	service := NewService(&repositoryStub{users: expected})

	actual, err := service.List(context.Background(), ListUsersQuery{Page: 1, Limit: 10})

	require.NoError(t, err)
	require.Equal(t, expected, actual.Users)
	require.Equal(t, 1, actual.Page)
	require.Equal(t, 1, actual.TotalPages)
}
