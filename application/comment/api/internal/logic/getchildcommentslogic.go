package logic

import (
	"context"
	"zhifou/application/comment/rpc/service"

	"zhifou/application/comment/api/internal/svc"
	"zhifou/application/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChildCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChildCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChildCommentsLogic {
	return &GetChildCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChildCommentsLogic) GetChildComments(req *types.GetChildCommentsRequest) (resp *types.GetChildCommentsResponse, err error) {
	rpcReq := &service.GetChildCommentsRequest{
		BizId:     req.BizId,
		ObjId:     req.ObjId,
		UserId:    req.UserId,
		ParentId:  req.ParentId,
		SortType:  req.SortType,
		Cursor:    req.Cursor,
		PageSize:  req.PageSize,
		LastObjId: req.LastObjId,
	}
	rpcResp, err := l.svcCtx.CommentRPC.GetChildComments(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to GetChildComments :", err)
		return nil, err
	}

	// 将 RPC 返回的 Items 转换为 API 层的类型
	items := make([]*types.CommentItem, len(rpcResp.Comments))
	for i, item := range rpcResp.Comments {
		items[i] = &types.CommentItem{
			CommentId:       item.CommentId,
			BizId:           item.BizId,
			ObjId:           item.ObjId,
			CommentUserId:   item.CommentUserId,
			BeCommentUserId: item.BeCommentUserId,
			ParentId:        item.ParentId,
			Content:         item.Content,
			LikeNum:         item.LikeNum,
			CreateTime:      item.CreateTime,
			CreateTimeUnix:  item.CreateTimeUnix,
		}
	}

	return &types.GetChildCommentsResponse{
		Comments: items,
		BizId:    rpcResp.BizId,
		IsEnd:    rpcResp.IsEnd,
		Cursor:   rpcResp.Cursor,
		LastId:   rpcResp.LastId,
	}, nil
}
