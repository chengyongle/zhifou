package logic

import (
	"context"
	"zhifou/application/article/api/internal/svc"
	"zhifou/application/article/api/internal/types"
	"zhifou/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeleteLogic {
	return &ArticleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDeleteLogic) ArticleDelete(req *types.ArticleDeleteRequest) (resp *types.ArticleDeleteResponse, err error) {
	rpcReq := &pb.ArticleDeleteRequest{
		UserId:    req.UserId,
		ArticleId: req.ArticleId,
	}
	_, err = l.svcCtx.ArticleRPC.ArticleDelete(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to ArticleDelete :", err)
		return nil, err
	}

	return &types.ArticleDeleteResponse{}, nil
}
