package logic

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"time"
	"zhifou/application/collect/mq/internal/model"
	"zhifou/application/collect/mq/internal/svc"
	"zhifou/application/collect/mq/internal/types"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type CollectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CollectLogic {
	return &CollectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CollectLogic) Consume(ctx context.Context, _, message string) error {
	logx.Infof("Consume msg val: %s", message)
	var msg *types.CollectMsg
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		logx.Errorf("Consume val: %s error: %v", message, err)
		return err
	}
	return l.CollectOperate(msg)
}

func (l *CollectLogic) CollectOperate(msg *types.CollectMsg) error {
	var err error
	switch msg.CollectType {
	case 1: //点赞
		// 事务 为两个表修改，保持数据一致性
		err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			//修改点赞关系表
			if msg.CollectRecordId == -1 { //如果没有记录
				err = model.NewCollectRecordModel(tx).Insert(l.ctx, &model.CollectRecord{
					BizID:         msg.BizId,
					ObjID:         msg.ObjId,
					UserID:        msg.UserId,
					CollectStatus: types.CollectStatuscollect,
					CreateTime:    time.Now(),
					UpdateTime:    time.Now(),
				})
			} else {
				err = model.NewCollectRecordModel(tx).UpdateFields(l.ctx, msg.CollectRecordId, map[string]interface{}{
					"collect_status": types.CollectStatuscollect,
				})
			}
			if err != nil {
				return err
			}
			//修改点赞计数表
			return model.NewCollectCountModel(tx).IncrCollectCount(l.ctx, msg.BizId, msg.ObjId)
		})
		if err != nil {
			l.Logger.Errorf("[Collect] Transaction error: %v", err)
			return err
		}
	case 2: //取消点赞
		err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			//修改点赞关系表
			err = model.NewCollectRecordModel(tx).UpdateFields(l.ctx, msg.CollectRecordId, map[string]interface{}{
				"collect_status": types.CollectStatusuncollect,
			})
			if err != nil {
				return err
			}
			//修改点赞计数表
			return model.NewCollectCountModel(tx).DecrCollectCount(l.ctx, msg.BizId, msg.ObjId)
		})
		if err != nil {
			l.Logger.Errorf("[Collect] Transaction error: %v", err)
			return err
		}
	}
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewCollectLogic(ctx, svcCtx)),
	}
}
