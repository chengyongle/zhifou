package logic

import (
	"context"
	"zhifou/application/article/api/internal/svc"
	"zhifou/application/article/api/internal/types"
	"zhifou/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchArticleLogic {
	return &SearchArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchArticleLogic) SearchArticle(req *types.SearchArticleRequest) (resp *types.SearchArticleResponse, err error) {
	rpcReq := &pb.SearchArticleRequest{
		Query:     req.Query,
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		Cursor:    req.Cursor,
		SortType:  req.SortType,
		ArticleId: req.ArticleId,
	}
	rpcResp, err := l.svcCtx.ArticleRPC.SearchArticle(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to SearchArticle :", err)
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
	return &types.SearchArticleResponse{
		Articles:      items,
		IsEnd:         rpcResp.IsEnd,
		PageNum:       rpcResp.PageNum,
		LastArticleId: rpcResp.LastArticleId,
		Cursor:        rpcResp.Cursor,
	}, nil
}
