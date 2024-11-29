package logic

import (
	"context"
	"zhifou/application/comment/rpc/service"

	"zhifou/application/comment/api/internal/svc"
	"zhifou/application/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentRequest) (resp *types.CreateCommentResponse, err error) {
	rpcReq := &service.CreateCommentRequest{
		BizId:           req.BizId,
		ObjId:           req.ObjId,
		CommentUserId:   req.CommentUserId,
		BeCommentUserId: req.BeCommentUserId,
		ParentId:        req.ParentId,
		Content:         req.Content,
	}
	rpcResp, err := l.svcCtx.CommentRPC.CreateComment(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to createComment :", err)
		return nil, err
	}

	return &types.CreateCommentResponse{
		CommentId: rpcResp.CommentId,
	}, nil
}
