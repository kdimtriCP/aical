package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
)

type UserInfo struct {
	Id    string `gorm:"column:id;type:varchar(255);primary_key;not null;" json:"id" form:"id" query:"id" validate:"required"`
	Name  string `gorm:"column:name;type:varchar(255);" json:"name" form:"name" query:"name"`
	Email string `gorm:"column:email;type:varchar(255);" json:"email" form:"email" query:"email"`
	Code  string `gorm:"column:code;type:varchar(255);" json:"code" form:"code" query:"code"`
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
func (u *UserRepo) CreateUser(ctx context.Context, code string) (string, error) {
	token, err := u.google.GetToken(ctx, code)
	if err != nil {
		return "", err
	}
	userInfo, err := u.google.GetUserInfo(ctx, token)
	if err != nil {
		return "", err
	}
	tx := u.data.db.Create(&UserInfo{
		Id:    userInfo.Id,
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Code:  code,
	})
	if tx.Error != nil {
		return "", tx.Error
	}
	return userInfo.Id, nil
}

// GetUserById .
func (u *UserRepo) GetUserById(ctx context.Context, id string) (*biz.UserInfo, error) {
	var userInfo UserInfo
	tx := u.data.db.Where("id = ?", id).First(&userInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &biz.UserInfo{
		Id:    userInfo.Id,
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Code:  userInfo.Code,
	}, nil
}

// GetUserByEmail .
func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (*biz.UserInfo, error) {
	var userInfo UserInfo
	tx := u.data.db.Where("email = ?", email).First(&userInfo)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &biz.UserInfo{
		Id:    userInfo.Id,
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Code:  userInfo.Code,
	}, nil
}

// ListUsers .
func (u *UserRepo) ListUsers(ctx context.Context) ([]*biz.UserInfo, error) {
	var userInfos []*UserInfo
	tx := u.data.db.Find(&userInfos)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var bizUserInfos []*biz.UserInfo
	for _, userInfo := range userInfos {
		bizUserInfos = append(bizUserInfos, &biz.UserInfo{
			Id:    userInfo.Id,
			Name:  userInfo.Name,
			Email: userInfo.Email,
			Code:  userInfo.Code,
		})
	}
	return bizUserInfos, nil
}
