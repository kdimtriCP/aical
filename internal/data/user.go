package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           string     `gorm:"column:id;type:varchar(255);primary_key;not null;" json:"id" form:"id" query:"id" validate:"required"`
	Name         string     `gorm:"column:name;type:varchar(255);" json:"name" form:"name" query:"name"`
	Email        string     `gorm:"column:email;type:varchar(255);" json:"email" form:"email" query:"email"`
	RefreshToken string     `gorm:"column:refresh_token;type:varchar(255);" json:"refresh_token" form:"refresh_token" query:"refresh_token"`
	Calendars    []Calendar `json:"calendar,omitempty"`
}

type UserRepo struct {
	data   *Data
	google *Google
	log    *log.Helper
}

func NewUserRepo(data *Data, ggl *Google, logger log.Logger) biz.UserRepo {
	return &UserRepo{data: data, google: ggl, log: log.NewHelper(logger)}
}

// CreateUser .
func (u *UserRepo) CreateUser(ctx context.Context, code string) (*biz.User, error) {
	u.log.Debugf("create user code: %s", code)
	token, err := u.google.GetToken(ctx, code)
	if err != nil {
		return nil, err
	}
	user := &User{}
	userInfo, err := u.google.GetUserInfo(ctx, token)
	userInfo.RefreshToken = token.RefreshToken
	if err != nil {
		return nil, err
	}
	tx := u.data.db.Where("id = ?", userInfo.ID).First(user)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	if tx.RowsAffected > 0 {
		if user.RefreshToken != token.RefreshToken {
			user.RefreshToken = token.RefreshToken
			tx = u.data.db.Save(user)
			if tx.Error != nil {
				return nil, tx.Error
			}
		}
		u.log.Infof("user already exists: %s", user.Name)
		return user.Biz(), nil
	}
	tx = u.data.db.Create(userInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return userInfo.Biz(), nil
}

// GetUserById .
func (u *UserRepo) GetUserById(ctx context.Context, id string) (*biz.User, error) {
	var userInfo *User
	tx := u.data.db.Where("id = ?", id).First(&userInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return userInfo.Biz(), nil
}

// GetUserByEmail .
func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	var userInfo *User
	tx := u.data.db.Where("email = ?", email).First(&userInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return userInfo.Biz(), nil
}

// ListUsers .
func (u *UserRepo) ListUsers(ctx context.Context) ([]*biz.User, error) {
	var userInfos []*User
	tx := u.data.db.Find(&userInfos)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var bizUserInfos []*biz.User
	for _, userInfo := range userInfos {
		bizUserInfos = append(bizUserInfos, userInfo.Biz())
	}
	return bizUserInfos, nil
}

// toBizUser .
func (user *User) Biz() *biz.User {
	if user == nil {
		return nil
	}
	return &biz.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
