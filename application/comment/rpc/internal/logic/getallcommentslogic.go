package logic

import (
	"cmp"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"slices"
	"strconv"
	"time"
	"zhifou/application/comment/rpc/internal/code"
	"zhifou/application/comment/rpc/internal/model"
	"zhifou/application/comment/rpc/internal/types"

	"zhifou/application/comment/rpc/internal/svc"
	"zhifou/application/comment/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllCommentsLogic {
	return &GetAllCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取帖子的全部评论列表
func (l *GetAllCommentsLogic) GetAllComments(in *service.GetAllCommentsRequest) (*service.GetAllCommentsResponse, error) {
	if in.BizId != types.ArticleBusiness {
		return nil, code.CommentBusinessTypeInvalid
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.ObjId <= 0 {
		return nil, code.ObjIdInvalid
	}
	if in.ParentId != 0 {
		return nil, code.ParentIdInvalid
	}
	if in.SortType != types.SortCreateTime && in.SortType != types.SortLikeCount {
		return nil, code.SortTypeInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		if in.SortType == types.SortCreateTime {
			in.Cursor = time.Now().Unix()
		} else {
			in.Cursor = types.DefaultSortLikeCursor
		}
	}
	var (
		sortField      string
		sortLikeNum    int64
		sortCreateTime string
	)
	if in.SortType == types.SortLikeCount {
		sortField = "like_num"
		sortLikeNum = in.Cursor
	} else {
		sortField = "create_time"
		sortCreateTime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
	}
	var (
		err                error
		isCache, isEnd     bool
		lastId, cursor, tu int64
		timestr            string
		curPage            []*service.CommentItem
		comments           []*model.Comment
	)
	commentIds, _ := l.cacheComments(l.ctx, in.BizId, in.ObjId, in.ParentId, in.Cursor, in.PageSize, in.SortType)
	if len(commentIds) > 0 {
		isCache = true
		if commentIds[len(commentIds)-1] == -1 {
			isEnd = true
		}
		comments, err = l.commentsByIds(l.ctx, commentIds)
		if err != nil {
			return nil, err
		}
		// 通过sortFiled对comments进行排序
		var cmpFunc func(a, b *model.Comment) int
		if sortField == "like_num" {
			cmpFunc = func(a, b *model.Comment) int {
				return cmp.Compare(b.LikeNum, a.LikeNum)
			}
		} else {
			cmpFunc = func(a, b *model.Comment) int {
				return cmp.Compare(b.CreateTime.Unix(), a.CreateTime.Unix())
			}
		}
		slices.SortFunc(comments, cmpFunc)
		//生成评论页面
		for _, comment := range comments {
			tu = comment.CreateTime.Unix()
			timestr = time.Unix(comment.CreateTime.Unix(), 0).Format("2006-01-02 15:04:05")
			curPage = append(curPage, &service.CommentItem{
				CommentId:       comment.Id,
				BizId:           comment.BizId,
				ObjId:           comment.ObjId,
				CommentUserId:   comment.CommentUserId,
				BeCommentUserId: comment.BeCommentUserId,
				ParentId:        comment.ParentId,
				Content:         comment.Content,
				LikeNum:         comment.LikeNum,
				CreateTime:      timestr,
				CreateTimeUnix:  tu,
			})
		}
	} else {
		//未命中缓存，数据库查找
		//singleflight应对缓存击穿
		v, err, _ := l.svcCtx.SingleFlightGroup.Do(fmt.Sprintf("CommentsByID:%d:%d:%d", in.BizId, in.ObjId, in.SortType), func() (interface{}, error) {
			return l.svcCtx.CommentModel.FindAllCommentsByBizIdAndObjId(l.ctx, in.BizId, in.ObjId, types.CommentStatusVisible, sortLikeNum, sortCreateTime, sortField, types.DefaultLimit)
		})
		if err != nil {
			logx.Errorf("FindAllCommentsByBizIdAndObjId BizId: %d ObjId: %d sortField: %s error: %v", in.BizId, in.ObjId, sortField, err)
			return nil, err
		}
		if v == nil {
			return &service.GetAllCommentsResponse{}, nil
		}
		comments = v.([]*model.Comment)
		var firstPageComments []*model.Comment
		if len(comments) > int(in.PageSize) { //如果返回的数量大于一页的数量，则切出要一页的数量来
			firstPageComments = comments[:int(in.PageSize)]
		} else { //如果返回的数量小于一页的数量，则说明不足一页，到最后了
			firstPageComments = comments
			isEnd = true
		}
		//生成评论页面
		for _, comment := range firstPageComments {
			tu = comment.CreateTime.Unix()
			timestr := time.Unix(comment.CreateTime.Unix(), 0).Format("2006-01-02 15:04:05")
			curPage = append(curPage, &service.CommentItem{
				CommentId:       comment.Id,
				BizId:           comment.BizId,
				ObjId:           comment.ObjId,
				CommentUserId:   comment.CommentUserId,
				BeCommentUserId: comment.BeCommentUserId,
				ParentId:        comment.ParentId,
				Content:         comment.Content,
				LikeNum:         comment.LikeNum,
				CreateTime:      timestr,
				CreateTimeUnix:  tu,
			})
		}
	}
	//出去重复条目（发布时间或点赞数相同）
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.CommentId
		if in.SortType == types.SortCreateTime {
			cursor = pageLast.CreateTimeUnix
		} else {
			cursor = int64(pageLast.LikeNum)
		}
		if cursor < 0 {
			cursor = 0
		}
		for k, comment := range curPage {
			if in.SortType == types.SortCreateTime {
				if comment.CreateTimeUnix == in.Cursor && comment.CommentId == in.LastObjId {
					curPage = curPage[k:]
					break
				}
			} else {
				if comment.LikeNum == in.Cursor && comment.CommentId == in.LastObjId {
					curPage = curPage[k:]
					break
				}
			}
		}
	}
	ret := &service.GetAllCommentsResponse{
		Comments: curPage,
		BizId:    in.BizId,
		IsEnd:    isEnd,
		Cursor:   cursor,
		LastId:   lastId,
	}
	if !isCache {
		threading.GoSafe(func() { //异步写入缓存
			if len(comments) < types.DefaultLimit && len(comments) > 0 { //最后一页，用 id=-1来表示
				comments = append(comments, &model.Comment{Id: -1})
			}
			err = l.addCacheComments(context.Background(), in.BizId, in.ObjId, in.ParentId, in.SortType, comments)
			if err != nil {
				logx.Errorf("addCacheComments error: %v", err)
			}
		})
	}

	return ret, nil
}

// 并行处理
func (l *GetAllCommentsLogic) commentsByIds(ctx context.Context, commentIds []int64) ([]*model.Comment, error) {
	comments, err := mr.MapReduce[int64, *model.Comment, []*model.Comment](func(source chan<- int64) {
		//生成数据
		for _, cid := range commentIds {
			if cid == -1 {
				continue
			}
			source <- cid
		}
	}, func(id int64, writer mr.Writer[*model.Comment], cancel func(error)) {
		//处理数据
		p, err := l.svcCtx.CommentModel.FindOne(ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(p)
	}, func(pipe <-chan *model.Comment, writer mr.Writer[[]*model.Comment], cancel func(error)) {
		//聚合数据返回
		var comments []*model.Comment
		for comment := range pipe {
			comments = append(comments, comment)
		}
		writer.Write(comments)
	})
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (l *GetAllCommentsLogic) cacheComments(ctx context.Context, bizId string, objId, parentId, cursor, pagesize int64, sortType int32) ([]int64, error) {
	key := getCommentsKey(bizId, types.AllComments, objId, parentId, sortType)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("ExistsCtx key: %s error: %v", key, err)
	}
	if b { //如果存在就使用Expire给缓存续期
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, types.CommentsExpireTime)
		if err != nil {
			logx.Errorf("ExpireCtx key: %s error: %v", key, err)
		}
	}
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pagesize))
	if err != nil {
		logx.Errorf("ZrevrangebyscoreWithScoresAndLimit key: %s error: %v", key, err)
		return nil, err
	}
	var ids []int64
	for _, pair := range pairs {
		id, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			logx.Errorf("strconv.ParseInt key: %s error: %v", pair.Key, err)
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (l *GetAllCommentsLogic) addCacheComments(ctx context.Context, bizId string, objId, parentId int64, sortType int32, comments []*model.Comment) error {
	if len(comments) == 0 {
		return nil
	}
	key := getCommentsKey(bizId, types.AllComments, objId, parentId, sortType)
	for _, comment := range comments {
		var score int64
		if comment.Id == -1 {
			score = 0
		} else if sortType == types.SortCreateTime {
			score = comment.CreateTime.Unix()
		} else {
			score = comment.LikeNum
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatInt(comment.Id, 10))
		if err != nil {
			logx.Errorf("[addCacheComment] BizRedis.ZaddCtx error: %v", err)
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, types.CommentsExpireTime)
}

func getCommentsKey(bizId, GetCommentsTypes string, objId, parentId int64, sortType int32) (keyStr string) {
	switch GetCommentsTypes {
	case types.AllComments:
		keyStr = fmt.Sprintf("biz#%s#%d#%s#%d", bizId, objId, types.AllComments, sortType)
	case types.RootComments:
		keyStr = fmt.Sprintf("biz#%s#%d#%s#%d", bizId, objId, types.RootComments, sortType)
	case types.ChildComments:
		keyStr = fmt.Sprintf("biz#%s#%d#%s#%d#%d", bizId, objId, types.ChildComments, parentId, sortType)
	}
	return keyStr
}
