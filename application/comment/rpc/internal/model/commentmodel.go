package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentModel = (*customCommentModel)(nil)
var cacheZhifouCommentPrefix = "cache:parentcommentdata:"

type (
	// CommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentModel.
	CommentModel interface {
		commentModel
		CommentTransactCtx(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		FindBizidAndObjIdAndParentIdAndBeCommentUserId(ctx context.Context, bizid string, objid, parentid, becommentuid int64) (*Comment, error)
		FindAllCommentsByBizIdAndObjId(ctx context.Context, bizId string, objId int64, status int, likeNum int64, pubTime, sortField string, limit int) ([]*Comment, error)
		FindRootCommentsByBizIdAndObjId(ctx context.Context, bizId string, objId int64, status int, likeNum int64, pubTime, sortField string, limit int) ([]*Comment, error)
		FindChildCommentsByBizIdAndObjId(ctx context.Context, bizId string, objId, parentId int64, status int, likeNum int64, pubTime, sortField string, limit int) ([]*Comment, error)
	}

	customCommentModel struct {
		*defaultCommentModel
	}
)

// NewCommentModel returns a model for the database table.
func NewCommentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn, c, opts...),
	}
}

func (m *defaultCommentModel) FindBizidAndObjIdAndParentIdAndBeCommentUserId(ctx context.Context, bizid string, objid, parentid, becommentuid int64) (*Comment, error) {
	// 构建缓存键
	zhifouCommentCacheKey := fmt.Sprintf("%s%s:%v:%v:%v", cacheZhifouCommentPrefix, bizid, objid, parentid, becommentuid)
	var resp Comment

	// 使用缓存查询，如果缓存中有数据，则直接返回
	err := m.QueryRowCtx(ctx, &resp, zhifouCommentCacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		// 构建 SQL 查询语句
		query := fmt.Sprintf(
			"select %s from %s where `biz_id` = ? and `obj_id` = ? and `id` = ? and `comment_user_id` = ? limit 1",
			commentRows, m.table)

		// 执行查询，将结果填充到 `v`（即 `resp`）中
		return conn.QueryRowCtx(ctx, v, query, bizid, objid, parentid, becommentuid)
	})

	// 错误处理
	switch err {
	case nil:
		// 如果查询成功，返回填充好的 Comment 对象
		return &resp, nil
	case sqlc.ErrNotFound:
		// 如果没有找到数据，返回自定义的 NotFound 错误
		return nil, ErrNotFound
	default:
		// 其他错误，返回原始错误
		return nil, err
	}
}
func (m *defaultCommentModel) CommentTransactCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})

}

func (m *defaultCommentModel) FindAllCommentsByBizIdAndObjId(ctx context.Context, bizId string, objId int64, status int, likeNum int64, createTime, sortField string, limit int) ([]*Comment, error) {
	var (
		err      error
		sql      string
		anyField any
		comments []*Comment
	)
	if sortField == "like_num" {
		anyField = likeNum
		sql = fmt.Sprintf("select "+commentRows+" from "+m.table+" where biz_id =? and obj_id=? and status=? and like_num < ? order by %s desc limit ?", sortField)
	} else {
		anyField = createTime
		sql = fmt.Sprintf("select "+commentRows+" from "+m.table+" where biz_id =? and obj_id=? and status=? and create_time < ? order by %s desc limit ?", sortField)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &comments, sql, bizId, objId, status, anyField, limit)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
func (m *defaultCommentModel) FindRootCommentsByBizIdAndObjId(ctx context.Context, bizId string, objId int64, status int, likeNum int64, createTime, sortField string, limit int) ([]*Comment, error) {
	var (
		err      error
		sql      string
		anyField any
		comments []*Comment
	)
	if sortField == "like_num" {
		anyField = likeNum
		sql = fmt.Sprintf("select "+commentRows+" from "+m.table+" where biz_id =? and obj_id=? and status=? and like_num < ? and parent_id = 0 order by %s desc limit ?", sortField)
	} else {
		anyField = createTime
		sql = fmt.Sprintf("select "+commentRows+" from "+m.table+" where biz_id =? and obj_id=? and status=? and create_time < ? and parent_id = 0 order by %s desc limit ?", sortField)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &comments, sql, bizId, objId, status, anyField, limit)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *defaultCommentModel) FindChildCommentsByBizIdAndObjId(ctx context.Context, bizId string, objId, parentId int64, status int, likeNum int64, createTime, sortField string, limit int) ([]*Comment, error) {
	var (
		err      error
		sql      string
		anyField any
		comments []*Comment
	)
	if sortField == "like_num" {
		anyField = likeNum
		sql = fmt.Sprintf("select "+commentRows+" from "+m.table+" where biz_id =? and obj_id=? and parent_id = ? and status=? and like_num < ?  order by %s desc limit ?", sortField)
	} else {
		anyField = createTime
		sql = fmt.Sprintf("select "+commentRows+" from "+m.table+" where biz_id =? and obj_id=? and parent_id = ? and status=? and create_time < ?  order by %s desc limit ?", sortField)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &comments, sql, bizId, objId, parentId, status, anyField, limit)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
