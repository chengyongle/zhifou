package logic

import (
	"context"
	"zhifou/application/collect/rpc/service"

	"zhifou/application/collect/api/internal/svc"
	"zhifou/application/collect/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnCollectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnCollectLogic {
	return &UnCollectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnCollectLogic) UnCollect(req *types.UnCollectRequest) (resp *types.UnCollectResponse, err error) {
	rpcReq := &service.UnCollectRequest{
		BizId:  req.BizId,
		ObjId:  req.ObjId,
		UserId: req.UserId,
	}
	rpcResp, err := l.svcCtx.CollectRPC.UnCollect(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to collect :", err)
		return nil, err
	}

	return &types.UnCollectResponse{
		BizId: rpcResp.BizId,
		ObjId: rpcResp.ObjId,
	}, nil
}
