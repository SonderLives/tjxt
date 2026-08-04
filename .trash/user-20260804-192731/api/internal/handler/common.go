package handler

import (
	"net/http"
	"strconv"

	"tjxt/pkg/response"
	"tjxt/pkg/xerr"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

// writeResult 将逻辑层响应写为标准 R 结构
func writeResult(w http.ResponseWriter, r *http.Request, resp *result.R, err error) {
	if err != nil {
		result.Write(w, r, nil, err)
		return
	}
	result.Write(w, r, resp.Data, nil)
}

// parsePathID 从路径参数解析 :id
func parsePathID(r *http.Request, req *types.UserIdReq) error {
	id, err := strconv.ParseInt(pathvar.Vars(r)["id"], 10, 64)
	if err != nil || id <= 0 {
		return xerr.New(xerr.CodeBadRequest, "路径参数 id 非法")
	}
	req.Id = id
	return nil
}