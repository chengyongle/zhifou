package logic

import (
	"context"
	"zhifou/application/comment/rpc/service"

	"zhifou/application/comment/api/internal/svc"
	"zhifou/application/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentRequest) (resp *types.DeleteCommentResponse, err error) {
	rpcReq := &service.DeleteCommentRequest{
		CommentId:     req.CommentId,
		CommentUserId: req.CommentUserId,
	}
	rpcResp, err := l.svcCtx.CommentRPC.DeleteComment(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to DeleteComment :", err)
		return nil, err
	}

	return &types.DeleteCommentResponse{
		CommentId: rpcResp.CommentId,
	}, nil
}
