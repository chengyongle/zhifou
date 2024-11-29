// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package handler

import (
	"net/http"

	"zhifou/application/comment/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/all",
				Handler: GetAllCommentsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/child",
				Handler: GetChildCommentsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: CreateCommentHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/delete",
				Handler: DeleteCommentHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/root",
				Handler: GetRootCommentsHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/v1/comment"),
	)
}
