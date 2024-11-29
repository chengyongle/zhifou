package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"zhifou/application/collect/rpc/internal/config"
	"zhifou/application/collect/rpc/internal/model"
	"zhifou/pkg/orm"
)

type ServiceContext struct {
	Config             config.Config
	KqPusherClient     *kq.Pusher
	DB                 *orm.DB
	CollectRecordModel *model.CollectRecordModel
	CollectCountModel  *model.CollectCountModel
	BizRedis           *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.DB.DataSource,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifetime:  c.DB.MaxLifetime,
	})
	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	return &ServiceContext{
		Config:             c,
		DB:                 db,
		KqPusherClient:     kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		CollectRecordModel: model.NewCollectRecordModel(db.DB),
		CollectCountModel:  model.NewCollectCountModel(db.DB),
		BizRedis:           rds,
	}
}
