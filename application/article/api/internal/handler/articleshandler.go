package handler

import (
	"net/http"
	"zhifou/application/article/api/internal/logic"
	"zhifou/application/article/api/internal/svc"
	"zhifou/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ArticlesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticlesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewArticlesLogic(r.Context(), svcCtx)
		resp, err := l.Articles(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
