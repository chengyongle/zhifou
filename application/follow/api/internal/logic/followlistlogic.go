package logic

import (
	"context"
	"zhifou/application/follow/rpc/pb"

	"zhifou/application/follow/api/internal/svc"
	"zhifou/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowListLogic) FollowList(req *types.FollowListRequest) (resp *types.FollowListResponse, err error) {
	rpcReq := &pb.FollowListRequest{
		Id:       req.Id,
		UserId:   req.UserId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	}
	rpcResp, err := l.svcCtx.FollowRPC.FollowList(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to collect :", err)
		return nil, err
	}

	// 将 RPC 返回的 Items 转换为 API 层的类型
	items := make([]*types.FollowItem, len(rpcResp.Items))
	for i, item := range rpcResp.Items {
		items[i] = &types.FollowItem{
			Id:             item.Id,
			FollowedUserId: item.FollowedUserId,
			FansCount:      item.FansCount,
			CreateTime:     item.CreateTime,
		}
	}

	return &types.FollowListResponse{
		Items:  items,
		Cursor: rpcResp.Cursor,
		IsEnd:  rpcResp.IsEnd,
		Id:     rpcResp.Id,
	}, nil
}
