package logic

import (
	"context"
	"zhifou/application/follow/rpc/pb"

	"zhifou/application/follow/api/internal/svc"
	"zhifou/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnFollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowLogic {
	return &UnFollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnFollowLogic) UnFollow(req *types.UnFollowRequest) (resp *types.UnFollowResponse, err error) {
	rpcReq := &pb.UnFollowRequest{
		UserId:         req.UserId,
		FollowedUserId: req.FollowedUserId,
	}
	_, err = l.svcCtx.FollowRPC.UnFollow(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to follow :", err)
		return nil, err
	}

	return &types.UnFollowResponse{}, nil
}
