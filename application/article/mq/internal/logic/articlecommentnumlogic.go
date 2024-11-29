package logic

import (
	"context"
	"encoding/json"
	"strconv"

	"zhifou/application/article/mq/internal/svc"
	"zhifou/application/article/mq/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleCommentNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleCommentNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleCommentNumLogic {
	return &ArticleCommentNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleCommentNumLogic) Consume(ctx context.Context, _, val string) error {
	logx.Infof("Consume msg val: %s", val)
	var msg *types.CanalCommentMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("Consume val: %s error: %v", val, err)
		return err
	}

	return l.updateArticleCommentNum(l.ctx, msg)
}

func (l *ArticleCommentNumLogic) updateArticleCommentNum(ctx context.Context, msg *types.CanalCommentMsg) error {
	if len(msg.Data) == 0 {
		return nil
	}

	for _, d := range msg.Data {
		if d.BizID != types.ArticleBizID {
			continue
		}
		objid, err := strconv.ParseInt(d.ObjID, 10, 64)
		if err != nil {
			logx.Errorf("strconv.ParseInt id: %s error: %v", d.ID, err)
			continue
		}
		commentNum, err := strconv.ParseInt(d.CommentNum, 10, 64)
		if err != nil {
			logx.Errorf("strconv.ParseInt commentNum: %s error: %v", d.CommentNum, err)
			continue
		}
		err = l.svcCtx.ArticleModel.UpdateCommentNum(ctx, objid, commentNum)
		if err != nil {
			logx.Errorf("UpdateCommentNum id: %d comment: %d", objid, commentNum)
		}
	}

	return nil
}
