package repository

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	minioClient "rip_project/internal/app/minioClient"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotAllowed    = errors.New("not allowed")
	ErrNoDraft       = errors.New("no draft for this user")
)

type Repository struct {
	db     *gorm.DB
	mc     *minio.Client
	redis  *redis.Client
	userID int
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	mc, err := minioClient.InitMinio()
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Repository{
		db:    db,
		mc:    mc,
		redis: redisClient,
	}, nil
}

func (r *Repository) GetRedis() *redis.Client {
	return r.redis
}

func (r *Repository) GetCreatorID() int {
	return 1
}

func (r *Repository) GetUserID() int {
	return r.userID
}
