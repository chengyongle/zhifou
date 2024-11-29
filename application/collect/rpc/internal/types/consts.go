package types

const (
	ArticleBusiness = "article" //文章类型
)
const (
	CollectStatuscollect   = iota + 1 // 收藏
	CollectStatusuncollect            // 取消收藏
)
const (
	DefaultPageSize      = 20
	CacheMaxCollectCount = 1000 // 缓存最大收藏数
)
