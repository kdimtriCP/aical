package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/internal/biz"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	GoogleID     string
	Name         string
	Email        string
	RefreshToken string
	Calendars    []*Calendar
}

// biz returns biz user.
func (u *User) biz() *biz.User {
	return &biz.User{
		ID:           u.ID,
		GoogleID:     u.GoogleID,
		Name:         u.Name,
		Email:        u.Email,
		RefreshToken: u.RefreshToken,
	}
}

// parseUser fills user from biz user.
func parseUser(bu *biz.User) *User {
	return &User{
		ID:           bu.ID,
		GoogleID:     bu.GoogleID,
		Name:         bu.Name,
		Email:        bu.Email,
		RefreshToken: bu.RefreshToken,
	}
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
func (r *UserRepo) Create(ctx context.Context, user *biz.User) error {
	r.log.Debugf("create u code: %v", user)
	u := parseUser(user)
	return r.data.db.Create(&u).Error
}

// Get gets user from database by id or email
func (r *UserRepo) Get(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.Debugf("get u: %v", user)
	u := parseUser(user)
	tx := r.data.db.Where(u).First(&u)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return u.biz(), nil
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
