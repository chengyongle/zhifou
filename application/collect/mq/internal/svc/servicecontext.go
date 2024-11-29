package svc

import (
	"zhifou/application/collect/mq/internal/config"
	"zhifou/application/collect/mq/internal/model"
	"zhifou/pkg/orm"
)

type ServiceContext struct {
	Config             config.Config
	DB                 *orm.DB
	CollectRecordModel *model.CollectRecordModel
	CollectCountModel  *model.CollectCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.DB.DataSource,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifetime:  c.DB.MaxLifetime,
	})
	return &ServiceContext{
		Config:             c,
		DB:                 db,
		CollectRecordModel: model.NewCollectRecordModel(db.DB),
		CollectCountModel:  model.NewCollectCountModel(db.DB),
	}
}
