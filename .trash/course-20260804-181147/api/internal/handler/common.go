package handler

import (
	"net/http"

	"tjxt/pkg/response"
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