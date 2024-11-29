package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"time"
	"zhifou/application/article/rpc/internal/code"
	"zhifou/application/article/rpc/internal/svc"
	"zhifou/application/article/rpc/internal/types"
	"zhifou/application/article/rpc/pb"
)

type SearchArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchArticleLogic {
	return &SearchArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchArticleLogic) SearchArticle(in *pb.SearchArticleRequest) (*pb.SearchArticleResponse, error) {
	if in.SortType != types.SortRelavance && in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount {
		return nil, code.SortTypeInvalid
	}
	if len(in.Query) == 0 {
		return nil, code.QueryContentInvalid
	}
	if in.PageNum < 1 {
		return nil, code.PageNumInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 && in.SortType != types.SortRelavance {
		if in.SortType == types.SortPublishTime {
			in.Cursor = time.Now().Unix()
		} else {
			in.Cursor = types.DefaultSortLikeCursor
		}
	}

	var (
		err            error
		isEnd          bool
		lastId, cursor int64
		curPage        []*pb.ArticleItem
	)
	// 执行搜索查询
	ESres, err := l.SearchArticlesInES(l.ctx, in.Query, in.PageNum, in.PageSize, in.Cursor, in.SortType)
	if err != nil {
		logx.Errorf("SearchArticlesInES query:%s,cursor: %d,pagesize: %d, error: %v", in.Query, in.Cursor, in.PageSize, err)
		return nil, err
	}
	if len(ESres) == 0 {
		return &pb.SearchArticleResponse{}, nil
	}
	//分页
	firstPageArticles := ESres
	if len(ESres) < int(in.PageSize) { //如果返回的数量大于一页的数量，则切出要一页的数量来
		isEnd = true
	}
	// 组装响应的文章列表
	for _, at := range firstPageArticles {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", at.PublishTime, time.Local)
		curPage = append(curPage, &pb.ArticleItem{
			Id:              at.ArticleId,
			Title:           at.Title,
			Content:         at.Content,
			LikeCount:       at.LikeNum,
			CollectCount:    at.CollectNum,
			CommentCount:    at.CommentNum,
			PublishTime:     at.PublishTime,
			PublishTimeUnix: t.Unix(),
		})
	}
	//出去重复条目（发布时间或点赞数相同）
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		if in.SortType == types.SortPublishTime {
			cursor = pageLast.PublishTimeUnix
		} else {
			cursor = pageLast.LikeCount
		}
		if cursor < 0 {
			cursor = 0
		}
		for k, article := range curPage {
			if in.SortType == types.SortPublishTime {
				if article.PublishTimeUnix == in.Cursor && article.Id == in.ArticleId {
					curPage = curPage[k:]
					break
				}
			} else {
				if article.LikeCount == in.Cursor && article.Id == in.ArticleId {
					curPage = curPage[k:]
					break
				}
			}
		}
	}
	ret := &pb.SearchArticleResponse{
		IsEnd:         isEnd,
		Cursor:        cursor,
		LastArticleId: lastId,
		PageNum:       in.PageNum,
		Articles:      curPage,
	}
	return ret, nil
}

// SearchArticlesInES  在 Elasticsearch 的文章索引中执行搜索查询
func (l *SearchArticleLogic) SearchArticlesInES(ctx context.Context, query string, pagenum, pagesize, cursor int64, sorttype int32) ([]*types.ArticleEsMsg, error) {
	// 创建搜索请求体
	var (
		from        int64
		buf         bytes.Buffer
		searchQuery map[string]interface{}
	)

	switch sorttype {
	case types.SortRelavance:

		from = (pagenum - 1) * pagesize
		searchQuery = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"multi_match": map[string]interface{}{
								"query":  query,                                       // 查询内容
								"fields": []string{"title", "content", "description"}, // 查询字段
							},
						},
					},
					"filter": []map[string]interface{}{
						{
							"term": map[string]interface{}{
								"status": 2, // 添加过滤条件，要求 status 字段值为 2
							},
						},
					},
				},
			},
			"from": from,
			"size": pagesize,
		}
	case types.SortLikeCount:
		searchQuery = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"multi_match": map[string]interface{}{
								"query":  query,                                       // 查询内容
								"fields": []string{"title", "content", "description"}, // 查询字段
							},
						},
					},
					"filter": []map[string]interface{}{
						{
							"range": map[string]interface{}{
								"like_num": map[string]interface{}{
									"lt": cursor, // 使用 cursor 变量作为 publish_time 的上限
								},
							},
						},
						{
							"term": map[string]interface{}{
								"status": 2, // 添加过滤条件，要求 status 字段值为 2
							},
						},
					},
				},
			},
			"from": 0,
			"size": pagesize,
			"sort": []map[string]interface{}{
				{
					"like_num": map[string]string{
						"order": "desc",
					},
				},
			},
		}
	case types.SortPublishTime:
		searchQuery = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"multi_match": map[string]interface{}{
								"query":  query,                                       // 查询内容
								"fields": []string{"title", "content", "description"}, // 查询字段
							},
						},
					},
					"filter": []map[string]interface{}{
						{
							"range": map[string]interface{}{
								"publish_time_unix": map[string]interface{}{
									"lt": cursor, // 使用 cursor 变量作为 publish_time 的上限
								},
							},
						},
						{
							"term": map[string]interface{}{
								"status": 2, // 添加过滤条件，要求 status 字段值为 2
							},
						},
					},
				},
			},
			"from": 0,
			"size": pagesize,
			"sort": []map[string]interface{}{
				{
					"publish_time_unix": map[string]string{
						"order": "desc",
					},
				},
			},
		}
	}
	// 将查询结构体编码为 JSON 并写入缓冲区
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		logx.Errorf("Error occurred while encoding search query: %v", err)
		return nil, err
	}

	//执行搜索请求
	res, err := l.svcCtx.Es.Client.Search(
		l.svcCtx.Es.Client.Search.WithContext(l.ctx),
		l.svcCtx.Es.Client.Search.WithIndex("article-index"), // 指定索引为 article-index
		l.svcCtx.Es.Client.Search.WithBody(&buf),
	)
	if err != nil {
		logx.Errorf("Error executing search request: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	//检查搜索错误
	err = l.Checkforsearcherrors(l.ctx, res)
	if err != nil {
		return nil, err
	}
	// 解析搜索结果
	var r struct {
		Hits struct {
			Hits []struct {
				Source *types.ArticleEsMsg `json:"_source"` // 将每个搜索命中的 _source 字段反序列化为 ArticleEsMsg 类型
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logx.Errorf("Error parsing search results: %v", err)
		return nil, err
	}

	// 从搜索结果中提取文章列表
	var articles []*types.ArticleEsMsg
	for _, hit := range r.Hits.Hits {
		articles = append(articles, hit.Source)
	}
	return articles, nil
}

// 检查搜索错误
func (l *SearchArticleLogic) Checkforsearcherrors(ctx context.Context, res *esapi.Response) error {
	// 读取并打印响应内容
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logx.Errorf("Error reading response body: %v", err)
		return err
	}
	defer res.Body.Close()

	// 打印原始响应内容，便于调试
	logx.Infof("Raw Response Body: %s", string(body))

	// 重新将内容写回 res.Body，以便后续解析
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// 解析 JSON 并检查 "error" 字段
	var esResponse map[string]interface{}
	if err := json.Unmarshal(body, &esResponse); err != nil {
		logx.Errorf("Error parsing JSON response: %v", err)
		return err
	}

	if esError, exists := esResponse["error"]; exists {
		// 如果存在 "error" 字段，说明请求出错，打印错误信息
		logx.Errorf("Elasticsearch  search error: %v", esError)
		return err
	}
	return nil
}
