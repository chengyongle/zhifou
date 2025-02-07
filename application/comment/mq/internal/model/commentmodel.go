package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentModel = (*customCommentModel)(nil)

type (
	// CommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentModel.
	CommentModel interface {
		commentModel
		withSession(session sqlx.Session) CommentModel
		UpdateLikeNum(ctx context.Context, objid, likeNum int64) error
	}

	customCommentModel struct {
		*defaultCommentModel
	}
)

// NewCommentModel returns a model for the database table.
func NewCommentModel(conn sqlx.SqlConn) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn),
	}
}

func (m *customCommentModel) withSession(session sqlx.Session) CommentModel {
	return NewCommentModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customCommentModel) UpdateLikeNum(ctx context.Context, objid, likeNum int64) error {
	query := fmt.Sprintf("update %s set like_num = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, likeNum, objid)
	return err
}
