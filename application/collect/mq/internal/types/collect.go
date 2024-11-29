package types

type CollectMsg struct {
	BizId           string ` json:"bizId,omitempty"`           // 业务id
	ObjId           int64  ` json:"objId,omitempty"`           // 收藏对象id
	UserId          int64  ` json:"userId,omitempty"`          // 用户id
	CollectType     int32  ` json:"collectType,omitempty"`     // 类型
	CollectRecordId int64  ` json:"collectRecordId,omitempty"` //存在的收藏记录id
}
