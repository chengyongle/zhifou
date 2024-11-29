package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"zhifou/application/comment/rpc/internal/code"
	"zhifou/application/comment/rpc/internal/svc"
	"zhifou/application/comment/rpc/internal/types"
	"zhifou/application/comment/rpc/service"
	"zhifou/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除评论
func (l *DeleteCommentLogic) DeleteComment(in *service.DeleteCommentRequest) (*service.DeleteCommentResponse, error) {
	if in.CommentId <= 0 {
		return nil, code.CommentIdInvalid
	}
	if in.CommentUserId <= 0 {
		return nil, code.UserIdInvalid
	}
	comment, err := l.svcCtx.CommentModel.FindOne(l.ctx, in.CommentId)
	if err != nil {
		l.Logger.Errorf("DeleteComment FindOne req: %v error: %v", in, err)
		return nil, err
	}
	if comment == nil || comment.Status == types.CommentStatusUserDelete {
		return &service.DeleteCommentResponse{}, nil
	}
	if comment.CommentUserId != in.CommentUserId {
		return nil, xcode.AccessDenied
	}
	//执行删除事务
	err = l.svcCtx.CommentModel.CommentTransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) (err error) {
		//修改comment表
		_, err = session.Exec(`UPDATE comment SET status = 1  WHERE id = ?;`, comment.Id)
		if err != nil {
			return err
		}
		//修改comment_count表
		if comment.ParentId == 0 { //如果是根评论
			_, err = session.Exec(`UPDATE comment_count SET comment_num=comment_num-1,comment_root_num = comment_root_num - 1  WHERE biz_id = ? AND obj_id = ? AND  comment_num > 0 AND  comment_root_num > 0;`, comment.BizId, comment.ObjId)

		} else {
			_, err = session.Exec(`UPDATE comment_count SET comment_num=comment_num-1 WHERE biz_id = ? AND obj_id = ? AND  comment_num > 0 ;`, comment.BizId, comment.ObjId)

		}
		return err

	})
	if err != nil {
		l.Logger.Errorf("[Comment] Transaction error: %v", err)
		return nil, err
	}
	//删除缓存
	//任意评论
	l.updateCommentCache(l.ctx, comment.BizId, types.AllComments, comment.Id, comment.ObjId, comment.ParentId)
	if comment.ParentId == 0 { //根评论
		l.updateCommentCache(l.ctx, comment.BizId, types.RootComments, comment.Id, comment.ObjId, comment.ParentId)
	}
	if comment.ParentId != 0 { //子评论
		l.updateCommentCache(l.ctx, comment.BizId, types.ChildComments, comment.Id, comment.ObjId, comment.ParentId)
	}
	return &service.DeleteCommentResponse{
		CommentId: comment.Id,
	}, nil
}

func (l *DeleteCommentLogic) updateCommentCache(ctx context.Context, bizId, GetCommentsTypes string, commentId, objId, parentId int64) {
	var (
		CreateTimeKey = getCommentsKey(bizId, GetCommentsTypes, objId, parentId, types.SortCreateTime)
		likeNumKey    = getCommentsKey(bizId, GetCommentsTypes, objId, parentId, types.SortLikeCount)
	)
	_, err := l.svcCtx.BizRedis.ZremCtx(l.ctx, CreateTimeKey, strconv.FormatInt(commentId, 10))
	if err != nil {
		l.Logger.Errorf("[DeleteCommentCache] error: %v", err)
		return
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, likeNumKey, strconv.FormatInt(commentId, 10))
	if err != nil {
		l.Logger.Errorf("[DeleteCommentCache] error: %v", err)
		return
	}
	return
}
