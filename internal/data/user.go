package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           string      `gorm:"column:id;type:varchar(255);primary_key;not null;" json:"id" form:"id" query:"id" validate:"required"`
	Name         string      `gorm:"column:name;type:varchar(255);" json:"name" form:"name" query:"name"`
	Email        string      `gorm:"column:email;type:varchar(255);" json:"email" form:"email" query:"email"`
	RefreshToken string      `gorm:"column:refresh_token;type:varchar(255);" json:"refresh_token" form:"refresh_token" query:"refresh_token"`
	Calendars    []*Calendar `gorm:"foreignKey:UserID;references:ID" json:"calendars,omitempty"`
}

// biz returns biz user.
func (u *User) biz() *biz.User {
	if u == nil {
		return nil
	}
	return &biz.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

// parse fills user from biz user.
func (u *User) parse(bu *biz.User) {
	u.ID = bu.ID
	u.Name = bu.Name
	u.Email = bu.Email
	u.RefreshToken = bu.RefreshToken
}

type Users []*User

// biz returns biz users
func (us Users) biz() []*biz.User {
	users := make([]*biz.User, len(us))
	for i, u := range us {
		users[i] = u.biz()
	}
	return users
}

type UserRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &UserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateUser .
func (r *UserRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.Debugf("create u code: %v", user)
	var u *User
	tx := r.data.db.Where("id = ?", user.ID).First(u)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	if tx.RowsAffected > 0 && u != nil {
		r.log.Infof("u already exists: %s", u.Name)
		return u.biz(), nil
	}
	u.parse(user)
	tx = r.data.db.Create(u)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return u.biz(), nil
}

// Get gets user from database by id or email
func (r *UserRepo) Get(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.Debugf("get u: %v", user)
	var u *User
	if user.ID != "" {
		tx := r.data.db.Where("id = ?", user.ID).First(&u)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return u.biz(), nil
	} else if user.Email != "" {
		tx := r.data.db.Where("email = ?", user.Email).First(&u)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return u.biz(), nil
	} else if user.Name != "" {
		tx := r.data.db.Where("name = ?", user.Name).First(&u)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return u.biz(), nil
	}
	r.log.Errorf("Get u bad request: %v", user)
	return user, nil
}

// List lists all users from database
func (r *UserRepo) List(ctx context.Context) ([]*biz.User, error) {
	var us *Users
	tx := r.data.db.Find(&us)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return us.biz(), nil
}
