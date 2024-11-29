package types

type ThumbupMsg struct {
	BizId        string ` json:"bizId,omitempty"`        // 业务id
	ObjId        int64  ` json:"objId,omitempty"`        // 点赞对象id
	UserId       int64  ` json:"userId,omitempty"`       // 用户id
	LikeType     int32  ` json:"likeType,omitempty"`     // 类型
	LikeRecordId int64  ` json:"likeRecordId,omitempty"` //存在的点赞记录id
}
