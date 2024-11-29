package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zhifou/application/collect/api/internal/logic"
	"zhifou/application/collect/api/internal/svc"
	"zhifou/application/collect/api/internal/types"
)

func CollectListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CollectListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCollectListLogic(r.Context(), svcCtx)
		resp, err := l.CollectList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
