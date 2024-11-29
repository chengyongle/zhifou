package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ArticleModel = (*customArticleModel)(nil)

type (
	// ArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleModel.
	ArticleModel interface {
		articleModel
		UpdateLikeNum(ctx context.Context, objid, likeNum int64) error
		UpdateCollectNum(ctx context.Context, objid int64, collectNum int64) error
		UpdateCommentNum(ctx context.Context, objid int64, commentNum int64) error
	}

	customArticleModel struct {
		*defaultArticleModel
	}
)

// NewArticleModel returns a model for the database table.
func NewArticleModel(conn sqlx.SqlConn) ArticleModel {
	return &customArticleModel{
		defaultArticleModel: newArticleModel(conn),
	}
}

func (m *customArticleModel) UpdateLikeNum(ctx context.Context, objid, likeNum int64) error {
	query := fmt.Sprintf("update %s set like_num = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, likeNum, objid)
	return err
}

func (m *customArticleModel) UpdateCollectNum(ctx context.Context, objid, collectNum int64) error {
	query := fmt.Sprintf("update %s set collect_num = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, collectNum, objid)
	return err
}

func (m *customArticleModel) UpdateCommentNum(ctx context.Context, objid, commentNum int64) error {
	query := fmt.Sprintf("update %s set comment_num = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, commentNum, objid)
	return err
}
