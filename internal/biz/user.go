package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	GoogleID     string    `json:"google_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, user *User) (*User, error)
	List(ctx context.Context) ([]*User, error)
}

type UserUseCase struct {
	db  UserRepo
	gg  GoogleRepo
	log *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		db:  repo,
		log: log.NewHelper(logger)}
}

// Create creates user in database
func (uc *UserUseCase) Create(ctx context.Context, user *User) error {
	uc.log.Debugf("Create user code: %v", user)
	return uc.db.Create(ctx, user)
}

// Get gets user from database
func (uc *UserUseCase) Get(ctx context.Context, user *User) (*User, error) {
	uc.log.Debugf("get user: %v", user)
	return uc.db.Get(ctx, user)
}

// ListUsers lists all users from database
func (uc *UserUseCase) List(ctx context.Context) ([]*User, error) {
	uc.log.Debugf("list users")
	return uc.db.List(ctx)
}
