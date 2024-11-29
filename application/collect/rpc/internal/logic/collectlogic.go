package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/threading"
	"strconv"
	"time"
	"zhifou/application/collect/rpc/internal/code"
	"zhifou/application/collect/rpc/internal/types"

	"zhifou/application/collect/rpc/internal/svc"
	"zhifou/application/collect/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type CollectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CollectLogic {
	return &CollectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 收藏
func (l *CollectLogic) Collect(in *service.CollectRequest) (*service.CollectResponse, error) {
	if in.BizId != types.ArticleBusiness {
		return nil, code.CollectBusinessTypeInvalid
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.ObjId <= 0 {
		return nil, code.ObjIdInvalid
	}
	collectrecordid := int64(-1)
	//查询是否收藏过
	iscollect, err := l.svcCtx.CollectRecordModel.FindByBizIDAndObjIDAndUserID(l.ctx, in.BizId, in.ObjId, in.UserId)
	if err != nil {
		l.Logger.Errorf("[Collect] CollectModel.FindByBizIDAndObjIDAndUserID err: %v req: %v", err, in)
		return nil, err
	}
	if iscollect != nil && iscollect.CollectStatus == types.CollectStatuscollect {
		return &service.CollectResponse{}, nil
	}
	if iscollect != nil {
		collectrecordid = iscollect.ID
	}
	//发给消息队列
	msg := &types.CollectMsg{
		BizId:           in.BizId,
		ObjId:           in.ObjId,
		UserId:          in.UserId,
		Collecttype:     1,
		CollectRecordId: collectrecordid,
	}
	// 发送kafka消息，异步
	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Logger.Errorf("[Collect] marshal msg: %v error: %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusherClient.Push(l.ctx, string(data))
		if err != nil {
			l.Logger.Errorf("[Collect] kq push data: %s error: %v", data, err)
		}

	})
	//修改缓存
	key := userCollectKey(in.BizId, in.UserId)
	userCollectExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, key)
	if err != nil {
		l.Logger.Errorf("[Follow] Redis Exists error: %v", err)
		return nil, err
	}
	if userCollectExist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, key, time.Now().Unix(), strconv.FormatInt(in.ObjId, 10))
		if err != nil {
			l.Logger.Errorf("[Collect] Redis Zadd error: %v", err)
			return nil, err
		}
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, key, 0, -(types.CacheMaxCollectCount + 1))
		if err != nil {
			l.Logger.Errorf("[Collect] Redis Zremrangebyrank error: %v", err)
		}
	}
	return &service.CollectResponse{
		BizId: in.BizId,
		ObjId: in.ObjId,
	}, nil
}
