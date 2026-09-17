package users

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

type ListUsersQuery struct {
	Page  int `form:"page" validate:"omitempty,min=1"`
	Limit int `form:"limit" validate:"omitempty,min=1,max=100"`
}

func (q *ListUsersQuery) Normalize() {
	if q.Page == 0 {
		q.Page = defaultPage
	}
	if q.Limit == 0 {
		q.Limit = defaultLimit
	}
	if q.Limit > maxLimit {
		q.Limit = maxLimit
	}
}
