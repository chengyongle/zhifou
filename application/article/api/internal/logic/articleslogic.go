package logic

import (
	"context"
	"zhifou/application/article/api/internal/svc"
	"zhifou/application/article/api/internal/types"
	"zhifou/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticlesLogic) Articles(req *types.ArticlesRequest) (resp *types.ArticlesResponse, err error) {
	rpcReq := &pb.ArticlesRequest{
		UserId:    req.UserId,
		Cursor:    req.Cursor,
		PageSize:  req.PageSize,
		SortType:  req.SortType,
		ArticleId: req.ArticleId,
	}
	rpcResp, err := l.svcCtx.ArticleRPC.Articles(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to Articles :", err)
		return nil, err
	}
	// 将 RPC 返回的 Items 转换为 API 层的类型
	items := make([]*types.ArticleItem, len(rpcResp.Articles))
	for i, item := range rpcResp.Articles {
		items[i] = &types.ArticleItem{
			Id:              item.Id,
			Title:           item.Title,
			Content:         item.Content,
			Description:     item.Description,
			Cover:           item.Cover,
			CommentCount:    item.CommentCount,
			LikeCount:       item.LikeCount,
			CollectCount:    item.CollectCount,
			PublishTimeUnix: item.PublishTimeUnix,
			PublishTime:     item.PublishTime,
			AuthorId:        item.AuthorId,
		}
	}
	return &types.ArticlesResponse{
		Items:         items,
		IsEnd:         rpcResp.IsEnd,
		Cursor:        rpcResp.Cursor,
		LastArticleId: rpcResp.LastArticleId,
	}, nil
}
