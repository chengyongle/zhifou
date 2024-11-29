package logic

import (
	"context"
	"zhifou/application/follow/rpc/pb"

	"zhifou/application/follow/api/internal/svc"
	"zhifou/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowLogic) Follow(req *types.FollowRequest) (resp *types.FollowResponse, err error) {
	rpcReq := &pb.FollowRequest{
		UserId:         req.UserId,
		FollowedUserId: req.FollowedUserId,
	}
	_, err = l.svcCtx.FollowRPC.Follow(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to follow :", err)
		return nil, err
	}

	return &types.FollowResponse{}, nil
}
