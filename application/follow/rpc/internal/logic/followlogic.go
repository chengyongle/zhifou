package logic

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
	"zhifou/application/follow/rpc/internal/code"
	"zhifou/application/follow/rpc/internal/types"
	"zhifou/application/follow/rpc/model"

	"zhifou/application/follow/rpc/internal/svc"
	"zhifou/application/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 关注
func (l *FollowLogic) Follow(in *pb.FollowRequest) (*pb.FollowResponse, error) {
	if in.UserId == 0 {
		return nil, code.FollowUserIdEmpty
	}
	if in.FollowedUserId == 0 {
		return nil, code.FollowedUserIdEmpty
	}
	if in.UserId == in.FollowedUserId {
		return nil, code.CannotFollowSelf
	}
	follow, err := l.svcCtx.FollowModel.FindByUserIDAndFollowedUserID(l.ctx, in.UserId, in.FollowedUserId)
	if err != nil {
		l.Logger.Errorf("[Follow] FollowModel.FindByUserIDAndFollowedUserID err: %v req: %v", err, in)
		return nil, err
	}
	if follow != nil && follow.FollowStatus == types.FollowStatusFollow {
		return &pb.FollowResponse{}, nil
	}
	// 事务 为两个表修改，保持数据一致性
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		//修改关注关系表
		if follow != nil {
			err = model.NewFollowModel(tx).UpdateFields(l.ctx, follow.ID, map[string]interface{}{
				"follow_status": types.FollowStatusFollow,
			})
		} else {
			err = model.NewFollowModel(tx).Insert(l.ctx, &model.Follow{
				UserID:         in.UserId,
				FollowedUserID: in.FollowedUserId,
				FollowStatus:   types.FollowStatusFollow,
				CreateTime:     time.Now(),
				UpdateTime:     time.Now(),
			})
		}
		if err != nil {
			return err
		}
		//修改关注计数表-用户的关注数
		err = model.NewFollowCountModel(tx).IncrFollowCount(l.ctx, in.UserId)
		if err != nil {
			return err
		}
		//修改关注计数表-被关注者的粉丝数
		return model.NewFollowCountModel(tx).IncrFansCount(l.ctx, in.FollowedUserId)
	})
	if err != nil {
		l.Logger.Errorf("[Follow] Transaction error: %v", err)
		return nil, err
	}
	//修改缓存
	followExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFollowKey(in.UserId))
	if err != nil {
		l.Logger.Errorf("[Follow] Redis Exists error: %v", err)
		return nil, err
	}
	if followExist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, userFollowKey(in.UserId), time.Now().Unix(), strconv.FormatInt(in.FollowedUserId, 10))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zadd error: %v", err)
			return nil, err
		}
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, userFollowKey(in.UserId), 0, -(types.CacheMaxFollowCount + 1))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zremrangebyrank error: %v", err)
		}
	}
	fansExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFansKey(in.FollowedUserId))
	if err != nil {
		l.Logger.Errorf("[Follow] Redis Exists error: %v", err)
		return nil, err
	}
	if fansExist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, userFansKey(in.FollowedUserId), time.Now().Unix(), strconv.FormatInt(in.UserId, 10))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zadd error: %v", err)
			return nil, err
		}
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, userFansKey(in.FollowedUserId), 0, -(types.CacheMaxFansCount + 1))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zremrangebyrank error: %v", err)
		}
	}
	return &pb.FollowResponse{}, nil
}

func userFollowKey(userId int64) string {
	return fmt.Sprintf("biz#user#follow#%d", userId)
}

func userFansKey(userId int64) string {
	return fmt.Sprintf("biz#user#fans#%d", userId)
}
