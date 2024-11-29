package logic

import (
	"context"
	"zhifou/application/article/api/internal/svc"
	"zhifou/application/article/api/internal/types"
	"zhifou/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDetailLogic) ArticleDetail(req *types.ArticleDetailRequest) (resp *types.ArticleDetailResponse, err error) {
	rpcReq := &pb.ArticleDetailRequest{
		ArticleId: req.ArticleId,
	}
	rpcResp, err := l.svcCtx.ArticleRPC.ArticleDetail(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to ArticleDetail :", err)
		return nil, err
	}
	item := &types.ArticleItem{
		Id:              rpcResp.Article.Id,
		Title:           rpcResp.Article.Title,
		Content:         rpcResp.Article.Content,
		Description:     rpcResp.Article.Description,
		Cover:           rpcResp.Article.Cover,
		CommentCount:    rpcResp.Article.CommentCount,
		LikeCount:       rpcResp.Article.LikeCount,
		CollectCount:    rpcResp.Article.CollectCount,
		PublishTimeUnix: rpcResp.Article.PublishTimeUnix,
		PublishTime:     rpcResp.Article.PublishTime,
		AuthorId:        rpcResp.Article.AuthorId,
	}
	return &types.ArticleDetailResponse{
		Article: item,
	}, nil
}
