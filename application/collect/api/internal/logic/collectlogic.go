package logic

import (
	"context"
	"zhifou/application/collect/rpc/service"

	"zhifou/application/collect/api/internal/svc"
	"zhifou/application/collect/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CollectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CollectLogic {
	return &CollectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CollectLogic) Collect(req *types.CollectRequest) (resp *types.CollectResponse, err error) {
	rpcReq := &service.CollectRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: req.UserId,
	}
	rpcResp, err := l.svcCtx.CollectRPC.Collect(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to collect :", err)
		return nil, err
	}

	return &types.CollectResponse{
		BizId: rpcResp.BizId,
		ObjId: rpcResp.ObjId,
	}, nil
}
