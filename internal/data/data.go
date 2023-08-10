package data

import (
	"github.com/kdimtricp/aical/internal/conf"
	"gorm.io/gorm/logger"
	slog "log"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewCache,
	NewAuthRepo,
	NewUserRepo,
	NewCalendarRepo,
	NewEventRepo,
	NewGoogleRepo,
	NewOpenAIRepo,
)

// Data .
type Data struct {
	db    *gorm.DB
	cache *redis.Client
}

// NewData .
func NewData(db *gorm.DB, cache *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Debug("closing the data resources")
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
			LogLevel:      logger.Warn,
		},
	)
	db, err = gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Errorf("failed opening connection to postgres: %v", err)
		return db, nil
	}
	tables := []interface{}{
		&User{},
		&Calendar{},
		&Event{},
		&EventHistory{},
	}
	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			log.Errorf("failed auto migrate table: %v", err)
			return db, nil
		}
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
		return nil, err
	}
	return cache, nil
}
