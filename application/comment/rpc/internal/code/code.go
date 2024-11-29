package code

import (
	"zhifou/pkg/xcode"
)

var (
	CommentBusinessTypeInvalid = xcode.New(90001, "评论业务类型无效")
	ObjIdInvalid               = xcode.New(90002, "对象ID无效")
	UserIdInvalid              = xcode.New(90003, "用户ID无效")
	BeCommentUserIdInvalid     = xcode.New(90004, "被评论用户ID无效")
	ParentIdInvalid            = xcode.New(90005, "父评论ID无效")
	CommentContentCantEmpty    = xcode.New(90006, "评论内容不能为空！")
	CommentIdInvalid           = xcode.New(90007, "评论ID无效！")
	SortTypeInvalid            = xcode.New(90008, "排序类型无效！")
)
