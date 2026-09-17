package users

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	List(ctx context.Context, offset, limit int) ([]User, int64, error)
}

type repository struct {
	db *gorm.DB
}

type userRecord struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey"`
	Name      string    `gorm:"column:name;not null"`
	Email     string    `gorm:"column:email;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (userRecord) TableName() string {
	return "users"
}

func newUserRecord(user User) userRecord {
	return userRecord{ID: user.ID, Name: user.Name, Email: user.Email, CreatedAt: user.CreatedAt}
}

func (record userRecord) user() User {
	return User{ID: record.ID, Name: record.Name, Email: record.Email, CreatedAt: record.CreatedAt}
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user User) (User, error) {
	record := newUserRecord(user)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return User{}, ErrEmailAlreadyExists
		}
		return User{}, err
	}
	return record.user(), nil
}

func (r *repository) List(ctx context.Context, offset, limit int) ([]User, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	records := make([]userRecord, 0)
	if err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	users := make([]User, 0, len(records))
	for _, record := range records {
		users = append(users, record.user())
	}
	return users, total, nil
}
