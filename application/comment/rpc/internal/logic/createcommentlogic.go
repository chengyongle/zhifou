package logic

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"time"
	"zhifou/application/comment/rpc/internal/code"
	"zhifou/application/comment/rpc/internal/types"

	"zhifou/application/comment/rpc/internal/svc"
	"zhifou/application/comment/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建评论
func (l *CreateCommentLogic) CreateComment(in *service.CreateCommentRequest) (*service.CreateCommentResponse, error) {
	if in.BizId != types.ArticleBusiness {
		return nil, code.CommentBusinessTypeInvalid
	}
	if in.ObjId <= 0 {
		return nil, code.ObjIdInvalid
	}
	if in.CommentUserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if len(in.Content) == 0 {
		return nil, code.CommentContentCantEmpty
	}
	//验证文章是否存在
	//at:=articlerpc(in.ObjId)
	if in.BeCommentUserId != 0 && in.ParentId == 0 {
		return nil, code.ParentIdInvalid
	}
	if in.ParentId != 0 {
		//验证父评论有效性
		_, err := l.svcCtx.CommentModel.FindBizidAndObjIdAndParentIdAndBeCommentUserId(l.ctx, in.BizId, in.ObjId, in.ParentId, in.BeCommentUserId)
		if err == sqlc.ErrNotFound {
			return nil, code.BeCommentUserIdInvalid
		}
		if err != nil {
			return nil, err
		}
	}
	//事务
	var res sql.Result
	err := l.svcCtx.CommentModel.CommentTransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) (err error) {
		//修改comment表
		res, err = session.Exec(`INSERT INTO comment (biz_id,obj_id,comment_user_id,be_comment_user_id,parent_id,content)VALUES (?,?,?,?,?,? )`, in.BizId, in.ObjId, in.CommentUserId, in.BeCommentUserId, in.ParentId, in.Content)
		if err != nil {
			return err
		}
		//修改comment_count表
		if in.ParentId == 0 { //如果是根评论
			_, err = session.Exec(`INSERT INTO comment_count (biz_id,obj_id, comment_num,comment_root_num) VALUES (?,?,1,1) ON DUPLICATE KEY UPDATE comment_num = comment_num + 1,comment_root_num = comment_root_num + 1;`, in.BizId, in.ObjId)

		} else {
			_, err = session.Exec(`INSERT INTO comment_count (biz_id,obj_id, comment_num,comment_root_num) VALUES (?,?,1,0) ON DUPLICATE KEY UPDATE comment_num = comment_num + 1;`, in.BizId, in.ObjId)

		}
		return err

	})
	if err != nil {
		l.Logger.Errorf("[Comment] Transaction error: %v", err)
		return nil, err
	}
	commentID, err := res.LastInsertId()
	if err != nil {
		l.Logger.Errorf("LastInsertId error: %v", err)
		return nil, err
	}
	//同步写缓存保证数据及时性
	//任意评论
	l.updateCommentCache(l.ctx, in.BizId, types.AllComments, commentID, in.ObjId, in.ParentId)
	if in.ParentId == 0 { //根评论
		l.updateCommentCache(l.ctx, in.BizId, types.RootComments, commentID, in.ObjId, in.ParentId)
	}
	if in.ParentId != 0 { //子评论
		l.updateCommentCache(l.ctx, in.BizId, types.ChildComments, commentID, in.ObjId, in.ParentId)
	}
	return &service.CreateCommentResponse{CommentId: commentID}, nil
}

func (l *CreateCommentLogic) updateCommentCache(ctx context.Context, bizId, GetCommentsTypes string, commentId, objId, parentId int64) {
	var (
		CreateTimeKey = getCommentsKey(bizId, GetCommentsTypes, objId, parentId, types.SortCreateTime)
		likeNumKey    = getCommentsKey(bizId, GetCommentsTypes, objId, parentId, types.SortLikeCount)
	)
	//同步写缓存保证数据及时性
	b, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, CreateTimeKey)
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, CreateTimeKey, time.Now().Unix(), strconv.FormatInt(commentId, 10))
		if err != nil {
			logx.Errorf("updateCommentCache  error: %v", err)
		}
	}
	if err != nil {
		logx.Errorf("l.svcCtx.BizRedis.ExistsCtx CreateTimeKey error: %v", err)
		return
	}
	b, err = l.svcCtx.BizRedis.ExistsCtx(l.ctx, likeNumKey)
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeNumKey, 0, strconv.FormatInt(commentId, 10))
		if err != nil {
			logx.Errorf("updateCommentCache  error: %v", err)
		}
	}
	if err != nil {
		logx.Errorf("l.svcCtx.BizRedis.ExistsCtx likeNumKey error: %v", err)
		return
	}
	return
}
