package logic

import (
	"context"
	"encoding/json"
	"strconv"

	"zhifou/application/article/mq/internal/svc"
	"zhifou/application/article/mq/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleCollectNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleCollectNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleCollectNumLogic {
	return &ArticleCollectNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleCollectNumLogic) Consume(ctx context.Context, _, val string) error {
	logx.Infof("Consume msg val: %s", val)
	var msg *types.CanalCollectMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("Consume val: %s error: %v", val, err)
		return err
	}

	return l.updateArticleCollectNum(l.ctx, msg)
}

func (l *ArticleCollectNumLogic) updateArticleCollectNum(ctx context.Context, msg *types.CanalCollectMsg) error {
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
		collectNum, err := strconv.ParseInt(d.CollectNum, 10, 64)
		if err != nil {
			logx.Errorf("strconv.ParseInt collectNum: %s error: %v", d.CollectNum, err)
			continue
		}
		err = l.svcCtx.ArticleModel.UpdateCollectNum(ctx, objid, collectNum)
		if err != nil {
			logx.Errorf("UpdateCollectNum id: %d collect: %d", objid, collectNum)
		}
	}

	return nil
}
