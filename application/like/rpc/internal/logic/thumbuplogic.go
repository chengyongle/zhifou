package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/threading"
	"zhifou/application/like/rpc/internal/code"
	"zhifou/application/like/rpc/internal/types"

	"zhifou/application/like/rpc/internal/svc"
	"zhifou/application/like/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *service.ThumbupRequest) (*service.ThumbupResponse, error) {
	if in.BizId != types.ArticleBusiness && in.BizId != types.CommentBusiness {
		return nil, code.LikeBusinessTypeInvalid
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.ObjId <= 0 {
		return nil, code.ObjIdInvalid
	}
	likerecordid := int64(-1)
	//查询是否点过赞
	islike, err := l.svcCtx.LikeRecordModel.FindByBizIDAndObjIDAndUserID(l.ctx, in.BizId, in.ObjId, in.UserId)
	if err != nil {
		l.Logger.Errorf("[Like] LikeModel.FindByBizIDAndObjIDAndUserID err: %v req: %v", err, in)
		return nil, err
	}
	if islike != nil && islike.LikeStatus == types.LikeStatuslike {
		return &service.ThumbupResponse{}, nil
	}
	if islike != nil {
		likerecordid = islike.ID
	}
	//发给消息队列
	msg := &types.ThumbupMsg{
		BizId:        in.BizId,
		ObjId:        in.ObjId,
		UserId:       in.UserId,
		Liketype:     1,
		LikeRecordId: likerecordid,
	}
	// 发送kafka消息，异步
	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Logger.Errorf("[Thumbup] marshal msg: %v error: %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusherClient.Push(l.ctx, string(data))
		if err != nil {
			l.Logger.Errorf("[Thumbup] kq push data: %s error: %v", data, err)
		}

	})
	return &service.ThumbupResponse{
		BizId: in.BizId,
		ObjId: in.ObjId,
	}, nil
}
