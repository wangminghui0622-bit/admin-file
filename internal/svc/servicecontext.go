package svc

import (
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-file/internal/cloud"
	"github.com/wangminghui0622-bit/admin-file/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	DB           *gorm.DB
	Query        *query.Query
	Redis        *redis.Redis
	CloudStorage *cloud.Group
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DatabaseDSN), &gorm.Config{})
	logx.Must(err)

	redisClient := redis.MustNewRedis(redis.RedisConf{
		Host: c.RedisConf.Host,
		Type: c.RedisConf.Type,
		Pass: c.RedisConf.Pass,
	})

	q := query.Use(db)

	return &ServiceContext{
		Config:       c,
		DB:           db,
		Query:        q,
		Redis:        redisClient,
		CloudStorage: cloud.NewGroup(db, q),
	}
}
