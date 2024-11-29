package logic

import (
	"context"
	"zhifou/application/like/rpc/service"

	"zhifou/application/like/api/internal/svc"
	"zhifou/application/like/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThumbupLogic) Thumbup(req *types.ThumbupRequest) (resp *types.ThumbupResponse, err error) {
	rpcReq := &service.ThumbupRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: req.UserId,
	}
	rpcResp, err := l.svcCtx.LikeRPC.Thumbup(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to Thumbup :", err)
		return nil, err
	}

	return &types.ThumbupResponse{
		BizId: rpcResp.BizId,
		ObjId: rpcResp.ObjId,
	}, nil
}
