package logic

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"time"
	"zhifou/application/like/mq/internal/model"
	"zhifou/application/like/mq/internal/types"

	"zhifou/application/like/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Consume(ctx context.Context, _, message string) error {
	logx.Infof("Consume msg val: %s", message)
	var msg *types.ThumbupMsg
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		logx.Errorf("Consume val: %s error: %v", message, err)
		return err
	}
	return l.ThumbupOperate(msg)
}

func (l *ThumbupLogic) ThumbupOperate(msg *types.ThumbupMsg) error {
	var err error
	//事务 为两个表修改，保持数据一致性
	switch msg.LikeType {
	case 1: //点赞
		// 事务 为两个表修改，保持数据一致性
		err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			//修改点赞关系表
			if msg.LikeRecordId == -1 { //如果没有记录
				err = model.NewLikeRecordModel(tx).Insert(l.ctx, &model.LikeRecord{
					BizID:      msg.BizId,
					ObjID:      msg.ObjId,
					UserID:     msg.UserId,
					LikeStatus: types.LikeStatuslike,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				})
			} else {
				err = model.NewLikeRecordModel(tx).UpdateFields(l.ctx, msg.LikeRecordId, map[string]interface{}{
					"like_status": types.LikeStatuslike,
				})
			}
			if err != nil {
				return err
			}
			//修改点赞计数表
			return model.NewLikeCountModel(tx).IncrLikeCount(l.ctx, msg.BizId, msg.ObjId)
		})
		if err != nil {
			l.Logger.Errorf("[Follow] Transaction error: %v", err)
			return err
		}
	case 2: //取消点赞
		err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			//修改点赞关系表
			err = model.NewLikeRecordModel(tx).UpdateFields(l.ctx, msg.LikeRecordId, map[string]interface{}{
				"like_status": types.LikeStatusunlike,
			})
			if err != nil {
				return err
			}
			//修改点赞计数表
			return model.NewLikeCountModel(tx).DecrLikeCount(l.ctx, msg.BizId, msg.ObjId)
		})
		if err != nil {
			l.Logger.Errorf("[Follow] Transaction error: %v", err)
			return err
		}
	}
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewThumbupLogic(ctx, svcCtx)),
	}
}
