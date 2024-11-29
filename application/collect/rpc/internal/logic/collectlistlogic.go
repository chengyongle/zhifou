package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/threading"
	"strconv"
	"time"
	"zhifou/application/collect/rpc/internal/code"
	"zhifou/application/collect/rpc/internal/model"
	"zhifou/application/collect/rpc/internal/svc"
	"zhifou/application/collect/rpc/internal/types"
	"zhifou/application/collect/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

const userCollectExpireTime = 3600 * 24 * 2

type CollectListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCollectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CollectListLogic {
	return &CollectListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 收藏列表
func (l *CollectListLogic) CollectList(in *service.CollectListRequest) (*service.CollectListResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = time.Now().Unix()
	}
	var (
		err                 error
		isCache, isEnd      bool
		lastId, cursor, oid int64
		t                   string
		collectrecords      []*model.CollectRecord
		curPage             []*service.CollectItem
	)
	p, _ := l.cacheCollectUserIds(l.ctx, in.BizId, in.UserId, in.Cursor, in.PageSize)
	if len(*p) > 0 {
		pairs := *p
		isCache = true
		lastItemId, _ := strconv.ParseInt(pairs[len(pairs)-1].Key, 10, 64)
		if lastItemId == -1 { //判断最后一页
			pairs = pairs[:len(pairs)-1]
			isEnd = true
		}
		if len(pairs) == 0 { //判空
			return &service.CollectListResponse{}, nil
		}
		//调用rpc查询详细的文章信息
		//collects, err = articlerpc()...

		for _, pair := range pairs {
			oid, _ = strconv.ParseInt(pair.Key, 10, 64)
			t = time.Unix(pair.Score, 0).Format("2006-01-02 15:04:05")
			curPage = append(curPage, &service.CollectItem{
				BizId:           in.BizId,
				ObjId:           oid,
				CollectTime:     t,
				CollectTimeUnix: pair.Score,
			})
		}
	} else {
		collectrecords, err = l.svcCtx.CollectRecordModel.FindByBizIDAndUserID(l.ctx, in.BizId, in.UserId, types.CacheMaxCollectCount)
		if err != nil {
			l.Logger.Errorf("[CollectList] CollectRecordModel.FindByBizIDAndUserID error: %v req: %v", err, in)
			return nil, err
		}
		if len(collectrecords) == 0 {
			return &service.CollectListResponse{}, nil
		}
		var firstPageCollectrecords []*model.CollectRecord //缓存中没有说明没访问过，这是必定是第一页
		if len(collectrecords) > int(in.PageSize) {
			firstPageCollectrecords = collectrecords[:in.PageSize]
		} else {
			firstPageCollectrecords = collectrecords
			isEnd = true
		}
		//调用rpc查询详细的文章信息
		//collects, err = articlerpc()...

		for _, collectrecord := range firstPageCollectrecords {
			curPage = append(curPage, &service.CollectItem{
				BizId:       in.BizId,
				ObjId:       collectrecord.ObjID,
				CollectTime: collectrecord.CreateTime.Format("2006-01-02 15:04:05"),
			})
		}
	}
	//去重
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.ObjId
		cursor = pageLast.CollectTimeUnix
		if cursor < 0 {
			cursor = 0
		}
		for k, collectItem := range curPage {
			if collectItem.CollectTimeUnix == in.Cursor && collectItem.ObjId == in.LastObjId {
				curPage = curPage[k:]
				break
			}
		}
	}

	//生成最后的展示结果
	ret := &service.CollectListResponse{
		BizId:  in.BizId,
		IsEnd:  isEnd,
		Cursor: cursor,
		LastId: lastId,
		Items:  curPage,
	}
	if !isCache {
		threading.GoSafe(func() {
			if len(collectrecords) < types.CacheMaxCollectCount && len(collectrecords) > 0 {
				collectrecords = append(collectrecords, &model.CollectRecord{ObjID: -1})
			}
			err = l.addCacheCollect(context.Background(), in.BizId, in.UserId, collectrecords)
			if err != nil {
				logx.Errorf("addCacheFollow error: %v", err)
			}
		})
	}
	return ret, nil
}
func (l *CollectListLogic) cacheCollectUserIds(ctx context.Context, bizId string, userId, cursor, pageSize int64) (*[]redis.Pair, error) {
	key := userCollectKey(bizId, userId)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("[cacheCollectObjIds] BizRedis.ExistsCtx error: %v", err)
	}
	if b {
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, userCollectExpireTime)
		if err != nil {
			logx.Errorf("[cacheCollectObjIds] BizRedis.ExpireCtx error: %v", err)
		}
	}
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	if err != nil {
		logx.Errorf("[cacheCollectObjIds] BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx error: %v", err)
		return nil, err
	}

	return &pairs, nil
}
func (l *CollectListLogic) addCacheCollect(ctx context.Context, bizId string, userId int64, collects []*model.CollectRecord) error {
	if len(collects) == 0 {
		return nil
	}
	key := userCollectKey(bizId, userId)
	for _, collect := range collects {
		var score int64
		if collect.ObjID == -1 {
			score = 0
		} else {
			score = collect.CreateTime.Unix()
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatInt(collect.ObjID, 10))
		if err != nil {
			logx.Errorf("[addCacheCollect] BizRedis.ZaddCtx error: %v", err)
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, userCollectExpireTime)
}

func userCollectKey(bizId string, userId int64) string {
	return fmt.Sprintf("biz#collect#%s#%d", bizId, userId)
}
