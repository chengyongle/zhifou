package types

// CommentLikeMsg canal解析like binlog消息.
type CommentLikeMsg struct {
	ID         string `json:"id"`
	BizID      string `json:"biz_id"`
	ObjID      string `json:"obj_id"`
	LikeNum    string `json:"like_num"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
