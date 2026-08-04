package handler

import (
	"errors"
	"net/http"
	"strconv"

	"common/result"
	"common/xerr"
	"user/internal/types"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

// writeResult 将逻辑层响应写为标准 R 结构：
// 成功时透出 data；失败时按业务错误码渲染。
func writeResult(w http.ResponseWriter, r *http.Request, resp *result.R, err error) {
	if err != nil {
		result.Write(w, r, nil, err)
		return
	}
	result.Write(w, r, resp.Data, nil)
}

// parsePathID 从路径参数解析 :id，避免先 Parse 消费请求体。
func parsePathID(r *http.Request, req *types.UserIdReq) error {
	vars := pathvar.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return xerr.New(xerr.CodeBadRequest, "缺少路径参数 id")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return errors.New("路径参数 id 非法")
	}
	req.Id = id
	return nil
}
