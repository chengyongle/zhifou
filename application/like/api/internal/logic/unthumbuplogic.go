package logic

import (
	"context"
	"zhifou/application/like/rpc/service"

	"zhifou/application/like/api/internal/svc"
	"zhifou/application/like/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnThumbupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnThumbupLogic {
	return &UnThumbupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnThumbupLogic) UnThumbup(req *types.UnThumbupRequest) (resp *types.UnThumbupResponse, err error) {
	rpcReq := &service.UnThumbupRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: req.UserId,
	}
	rpcResp, err := l.svcCtx.LikeRPC.UnThumbup(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to UnThumbup :", err)
		return nil, err
	}

	return &types.UnThumbupResponse{
		BizId: rpcResp.BizId,
		ObjId: rpcResp.ObjId,
	}, nil
}
