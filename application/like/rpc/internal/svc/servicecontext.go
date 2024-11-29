package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"zhifou/application/like/rpc/internal/config"
	"zhifou/application/like/rpc/internal/model"
	"zhifou/pkg/orm"
)

type ServiceContext struct {
	Config          config.Config
	KqPusherClient  *kq.Pusher
	DB              *orm.DB
	LikeRecordModel *model.LikeRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.DB.DataSource,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifetime:  c.DB.MaxLifetime,
	})
	return &ServiceContext{
		Config:          c,
		KqPusherClient:  kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		DB:              db,
		LikeRecordModel: model.NewLikeRecordModel(db.DB),
	}
}
