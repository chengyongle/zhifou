package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"zhifou/application/comment/mq/internal/svc"
	"zhifou/application/comment/mq/internal/types"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type CommentLikeNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentLikeNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLikeNumLogic {
	return &CommentLikeNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentLikeNumLogic) Consume(ctx context.Context, _, val string) error {
	logx.Infof("Consume msg val: %s", val)
	var msg *types.CommentLikeMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("Consume val: %s error: %v", val, err)
		return err
	}
	return l.updateCommentLikeNum(l.ctx, msg)
}

func (l *CommentLikeNumLogic) updateCommentLikeNum(ctx context.Context, msg *types.CommentLikeMsg) (err error) {
	if msg.BizID != types.CommentBizID {
		logx.Errorf("msg.BizID != types.CommentBizID")
		return nil
	}
	objid, err := strconv.ParseInt(msg.ObjID, 10, 64)
	if err != nil {
		logx.Errorf("strconv.ParseInt id: %s error: %v", msg.ID, err)
		return err
	}
	likeNum, err := strconv.ParseInt(msg.LikeNum, 10, 64)
	if err != nil {
		logx.Errorf("strconv.ParseInt likeNum: %s error: %v", msg.LikeNum, err)
		return err
	}
	err = l.svcCtx.CommentModel.UpdateLikeNum(ctx, objid, likeNum)
	if err != nil {
		logx.Errorf("UpdateLikeNum id: %d like: %d", objid, likeNum)
	}

	return err
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.CommentLikeKqConsumerConf, NewCommentLikeNumLogic(ctx, svcCtx)),
	}
}
