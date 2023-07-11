package data

import (
	"github.com/kdimtricp/aical/internal/conf"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm/logger"
	slog "log"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"golang.org/x/oauth2"
	calendarAPI "google.golang.org/api/calendar/v3"
	oauth2API "google.golang.org/api/oauth2/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewCache,
	NewGoogle,
	NewOpenAI,
	NewAuthRepo,
)

// Data .
type Data struct {
	db    *gorm.DB
	cache *redis.Client
}

// Google .
type Google struct {
	config *oauth2.Config
}

// OpenAI .
type OpenAI struct {
}

// NewData .
func NewData(c *conf.Data, db *gorm.DB, cache *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:    db,
		cache: cache,
	}, cleanup, nil
}

// NewDB .
func NewDB(c *conf.Data) (db *gorm.DB, err error) {
	newLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags|slog.Lshortfile),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)
	db, err = gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Errorf("failed opening connection to postgres: %v", err)
		return db, nil
	}
	return db, nil
}

// NewCache .
func NewCache(c *conf.Data) (cache *redis.Client, err error) {
	cache = redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       0,
	})
	_, err = cache.Ping().Result()
	if err != nil {
		log.Errorf("failed opening connection to redis: %v", err)
		return cache, nil
	}
	return cache, nil
}

// NewGoogle .
func NewGoogle(c *conf.Google, logger log.Logger) (*Google, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the google resources")
	}
	return &Google{
		config: &oauth2.Config{
			ClientID:     c.Client.Id,
			ClientSecret: c.Client.Secret,
			RedirectURL:  c.RedirectUrl,
			Scopes: []string{
				calendarAPI.CalendarScope,
				calendarAPI.CalendarEventsScope,
				oauth2API.UserinfoEmailScope,
				oauth2API.UserinfoProfileScope,
			},
			Endpoint: google.Endpoint,
		},
	}, cleanup, nil
}

// NewOpenAI .
func NewOpenAI(c *conf.OpenAI, logger log.Logger) (*OpenAI, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the openai resources")
	}
	return &OpenAI{}, cleanup, nil
}
