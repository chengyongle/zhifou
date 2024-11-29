package types

const (
	ArticleBusiness = "article" //文章类型
)
const (
	CommentStatusVisible    = iota // 评论状态可见
	CommentStatusUserDelete        // 用户删除
)
const (
	SortCreateTime = iota
	SortLikeCount
)

const (
	DefaultPageSize       = 20
	DefaultLimit          = 200
	DefaultSortLikeCursor = 1 << 30
)

const (
	AllComments   = "allcomments"
	RootComments  = "rootcomments"
	ChildComments = "childcomments"
)
const CommentsExpireTime = 3600 * 24 * 2
