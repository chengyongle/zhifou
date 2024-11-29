package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zhifou/application/comment/api/internal/logic"
	"zhifou/application/comment/api/internal/svc"
	"zhifou/application/comment/api/internal/types"
)

func GetAllCommentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAllCommentsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetAllCommentsLogic(r.Context(), svcCtx)
		resp, err := l.GetAllComments(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
