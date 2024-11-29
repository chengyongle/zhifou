package logic

import (
	"context"
	"zhifou/application/collect/rpc/service"

	"zhifou/application/collect/api/internal/svc"
	"zhifou/application/collect/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CollectListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCollectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CollectListLogic {
	return &CollectListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CollectListLogic) CollectList(req *types.CollectListRequest) (resp *types.CollectListResponse, err error) {
	rpcReq := &service.CollectListRequest{
		BizId:     req.BizId,
		UserId:    req.UserId,
		Cursor:    req.Cursor,
		PageSize:  req.PageSize,
		LastObjId: req.LastObjId,
	}
	rpcResp, err := l.svcCtx.CollectRPC.CollectList(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to collect :", err)
		return nil, err
	}

	// 将 RPC 返回的 Items 转换为 API 层的类型
	items := make([]*types.CollectItem, len(rpcResp.Items))
	for i, item := range rpcResp.Items {
		items[i] = &types.CollectItem{
			BizId:           item.BizId,
			ObjId:           item.ObjId,
			CollectTime:     item.CollectTime,
			CollectTimeUnix: item.CollectTimeUnix,
		}
	}

	return &types.CollectListResponse{
		Items:  items,
		BizId:  rpcResp.BizId,
		Cursor: rpcResp.Cursor,
		IsEnd:  rpcResp.IsEnd,
		LastId: rpcResp.LastId,
	}, nil
}
