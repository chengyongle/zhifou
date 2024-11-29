package code

import (
	"zhifou/pkg/xcode"
)

var (
	SortTypeInvalid         = xcode.New(60001, "排序类型无效")   // 排序类型无效
	UserIdInvalid           = xcode.New(60002, "用户ID无效")   // 用户ID无效
	ArticleTitleCantEmpty   = xcode.New(60003, "文章标题不能为空") // 文章标题不能为空
	ArticleContentCantEmpty = xcode.New(60004, "文章内容不能为空") // 文章内容不能为空
	ArticleIdInvalid        = xcode.New(60005, "文章ID无效")   // 文章ID无效
	QueryContentInvalid     = xcode.New(60006, "搜索内容不能为空") // 搜索内容不能为空
	PageNumInvalid          = xcode.New(60007, "页数必须为正整数") // 页数必须为正整数

)
