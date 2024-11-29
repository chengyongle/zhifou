package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
	"zhifou/application/article/mq/internal/config"
	"zhifou/application/article/mq/internal/model"
	"zhifou/application/user/rpc/user"
	"zhifou/pkg/es"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                    config.Config
	CommentLikeKqPusherClient *kq.Pusher
	ArticleModel              model.ArticleModel
	BizRedis                  *redis.Redis
	UserRPC                   user.User
	Es                        *es.Es
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

	conn := sqlx.NewMysql(c.Datasource)
	return &ServiceContext{
		Config:                    c,
		CommentLikeKqPusherClient: kq.NewPusher(c.CommentLikeKqPusherConf.Brokers, c.CommentLikeKqPusherConf.Topic),
		ArticleModel:              model.NewArticleModel(conn),
		BizRedis:                  rds,
		UserRPC:                   user.NewUser(zrpc.MustNewClient(c.UserRPC)),
		Es: es.MustNewEs(&es.Config{
			Addresses: c.Es.Addresses,
			Username:  c.Es.Username,
			Password:  c.Es.Password,
		}),
	}
}
