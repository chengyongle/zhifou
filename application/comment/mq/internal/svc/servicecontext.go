package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zhifou/application/comment/mq/internal/config"
	"zhifou/application/comment/mq/internal/model"
)

type ServiceContext struct {
	Config       config.Config
	CommentModel model.CommentModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Datasource)
	return &ServiceContext{
		Config:       c,
		CommentModel: model.NewCommentModel(conn),
	}
}
