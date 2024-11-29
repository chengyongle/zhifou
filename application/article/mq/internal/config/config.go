package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf

	LikeKqConsumerConf      kq.KqConf
	ArticleKqConsumerConf   kq.KqConf
	CollectKqConsumerConf   kq.KqConf
	CommentKqConsumerConf   kq.KqConf
	CommentLikeKqPusherConf struct {
		Brokers []string
		Topic   string
	}
	Datasource string
	BizRedis   redis.RedisConf
	// es config
	Es struct {
		Addresses []string
		Username  string
		Password  string
	}
	UserRPC zrpc.RpcClientConf
}
