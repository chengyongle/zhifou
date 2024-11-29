package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/singleflight"
	"zhifou/application/comment/rpc/internal/config"
	"zhifou/application/comment/rpc/internal/model"
)

type ServiceContext struct {
	Config            config.Config
	CommentModel      model.CommentModel
	CommentCountModel model.CommentCountModel
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:            c,
		CommentModel:      model.NewCommentModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		CommentCountModel: model.NewCommentCountModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:          rds,
	}
}
