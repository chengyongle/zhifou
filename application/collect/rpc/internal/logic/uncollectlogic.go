package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/threading"
	"strconv"
	"zhifou/application/collect/rpc/internal/code"
	"zhifou/application/collect/rpc/internal/types"

	"zhifou/application/collect/rpc/internal/svc"
	"zhifou/application/collect/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnCollectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnCollectLogic {
	return &UnCollectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消收藏
func (l *UnCollectLogic) UnCollect(in *service.UnCollectRequest) (*service.UnCollectResponse, error) {
	if in.BizId != types.ArticleBusiness {
		return nil, code.CollectBusinessTypeInvalid
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.ObjId <= 0 {
		return nil, code.ObjIdInvalid
	}
	//查询是否点过赞
	islike, err := l.svcCtx.CollectRecordModel.FindByBizIDAndObjIDAndUserID(l.ctx, in.BizId, in.ObjId, in.UserId)
	if err != nil {
		l.Logger.Errorf("[Like] LikeModel.FindByBizIDAndObjIDAndUserID err: %v req: %v", err, in)
		return nil, err
	}
	if islike == nil || islike.CollectStatus == types.CollectStatusuncollect {
		return &service.UnCollectResponse{}, nil
	}
	likerecordid := islike.ID
	//发给消息队列
	msg := &types.CollectMsg{
		BizId:           in.BizId,
		ObjId:           in.ObjId,
		UserId:          in.UserId,
		Collecttype:     2,
		CollectRecordId: likerecordid,
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
	//删除缓存
	key := userCollectKey(in.BizId, in.UserId)
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, key, strconv.FormatInt(in.ObjId, 10))
	if err != nil {
		l.Logger.Errorf("[UnCollect] BizRedis.ZremCtx error: %v", err)
		return nil, err
	}
	return &service.UnCollectResponse{
		BizId: in.BizId,
		ObjId: in.ObjId,
	}, nil
}
