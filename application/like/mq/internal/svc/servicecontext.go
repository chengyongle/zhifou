package svc

import (
	"zhifou/application/like/mq/internal/config"
	"zhifou/application/like/mq/internal/model"
	"zhifou/pkg/orm"
)

type ServiceContext struct {
	Config          config.Config
	DB              *orm.DB
	LikeRecordModel *model.LikeRecordModel
	LikeCountModel  *model.LikeCountModel
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
		DB:              db,
		LikeRecordModel: model.NewLikeRecordModel(db.DB),
		LikeCountModel:  model.NewLikeCountModel(db.DB),
	}
}
