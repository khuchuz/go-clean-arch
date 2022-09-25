package domain

import (
	"context"
	"time"
)

// User ...
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// UserUsecase represent the user's usecases
type UserUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, ar *User) error
	GetByEmail(ctx context.Context, email string) (User, error)
	Signup(context.Context, *User) error
	Delete(ctx context.Context, id int64) error
}

// UserRepository represent the user's repository contract
type UserRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []User, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, ar *User) error
	Signup(ctx context.Context, a *User) error
	Delete(ctx context.Context, id int64) error
}
