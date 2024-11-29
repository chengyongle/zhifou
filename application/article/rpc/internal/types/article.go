package types

type ArticleEsMsg struct {
	ArticleId   int64   `json:"article_id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Description string  `json:"description"`
	AuthorId    int64   `json:"author_id"`
	AuthorName  string  `json:"author_name"`
	Status      int     `json:"status"`
	CommentNum  int64   `json:"comment_num"`
	LikeNum     int64   `json:"like_num"`
	CollectNum  int64   `json:"collect_num"`
	ViewNum     int64   `json:"view_num"`
	ShareNum    int64   `json:"share_num"`
	TagIds      []int64 `json:"tag_ids"`
	PublishTime string  `json:"publish_time"`
	CreateTime  string  `json:"create_time"`
	UpdateTime  string  `json:"update_time"`
}
