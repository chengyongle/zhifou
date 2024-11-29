package types

type CollectMsg struct {
	BizId           string ` json:"bizId,omitempty"`           // 业务id
	ObjId           int64  ` json:"objId,omitempty"`           // 收藏对象id
	UserId          int64  ` json:"userId,omitempty"`          // 用户id
	Collecttype     int    ` json:"collectType,omitempty"`     // 收藏类型，收藏1还是取消2
	CollectRecordId int64  ` json:"collectRecordId,omitempty"` //存在的收藏记录id
}
