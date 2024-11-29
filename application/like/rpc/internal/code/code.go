package code

import (
	"zhifou/pkg/xcode"
)

var (
	LikeBusinessTypeInvalid = xcode.New(70001, "点赞类型无效") // 点赞业务类型无效
	UserIdInvalid           = xcode.New(70002, "用户ID无效") // 用户ID无效
	ObjIdInvalid            = xcode.New(70003, "用户ID无效") // 点赞对象ID无效
)
