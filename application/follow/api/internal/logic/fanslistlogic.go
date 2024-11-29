package logic

import (
	"context"
	"zhifou/application/follow/rpc/pb"

	"zhifou/application/follow/api/internal/svc"
	"zhifou/application/follow/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FansListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FansListLogic {
	return &FansListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FansListLogic) FansList(req *types.FansListRequest) (resp *types.FansListResponse, err error) {
	rpcReq := &pb.FansListRequest{
		UserId:   req.UserId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	}
	rpcResp, err := l.svcCtx.FollowRPC.FansList(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to collect :", err)
		return nil, err
	}

	// 将 RPC 返回的 Items 转换为 API 层的类型
	items := make([]*types.FansItem, len(rpcResp.Items))
	for i, item := range rpcResp.Items {
		items[i] = &types.FansItem{
			UserId:     item.UserId,
			FansUserId: item.FansUserId,
			CreateTime: item.CreateTime,
		}
	}

	return &types.FansListResponse{
		Items:  items,
		Cursor: rpcResp.Cursor,
	}, nil
}
