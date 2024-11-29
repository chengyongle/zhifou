package code

import (
	"zhifou/pkg/xcode"
)

var (
	CollectBusinessTypeInvalid = xcode.New(80001, "收藏类型无效")   // 收藏业务类型无效
	UserIdInvalid              = xcode.New(80002, "用户ID无效")   // 用户ID无效
	ObjIdInvalid               = xcode.New(80003, "收藏对象ID无效") // 收藏对象ID无效
)
