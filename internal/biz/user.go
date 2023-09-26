package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	GoogleID     string    `json:"google_id"`
	TGID         string    `json:"tgid"`
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

// GetUserByID gets user from database by ID string
func (uc *UserUseCase) GetUserByID(ctx context.Context, id string) (*User, error) {
	uc.log.Debugf("get user by ID: %v", id)
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID: uid,
	}
	return uc.db.Get(ctx, user)
}

// GetUserByTGID gets user from database by TGID string
func (uc *UserUseCase) GetUserByTGID(ctx context.Context, tgid string) (*User, error) {
	uc.log.Debugf("get user by TGID: %v", tgid)
	user := &User{
		TGID: tgid,
	}
	return uc.db.Get(ctx, user)
}
