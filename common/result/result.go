// Package result 提供统一的 API 响应结构 R。
//
// 接口契约见 Apifox 中统一响应模型 R{code,msg,requestId,data}。
// 所有服务通过 result.Write 输出标准响应，请求链路 ID 取自 trace。
package result

import (
	"net/http"

	"common/xerr"

	"github.com/zeromicro/go-zero/core/trace"
)

// R 统一响应体。
type R struct {
	Code      int64  `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"requestId"`
	Data      any    `json:"data"`
}

// Page 通用分页数据结构，对应 Apifox 中的 PageDTO。
type Page struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Pages int64 `json:"pages"`
}

// Ok 返回成功响应。
func Ok() *R {
	return &R{Code: int64(xerr.CodeSuccess), Msg: xerr.MsgSuccess, RequestId: "", Data: nil}
}

// OkData 返回携带数据的成功响应。
func OkData(data any) *R {
	return &R{Code: int64(xerr.CodeSuccess), Msg: xerr.MsgSuccess, RequestId: "", Data: data}
}

// Fail 由 error 构造失败响应，提取业务错误码与提示信息。
func Fail(err error) *R {
	if err == nil {
		return Ok()
	}
	return &R{
		Code: int64(xerr.CodeOf(err)),
		Msg:  xerr.MsgOf(err),
		Data: nil,
	}
}

// Write 将 data / err 渲染为标准响应写入 w。
//
// - err == nil 时写入 code=200 的成功响应；
// - err 为业务错误时写入对应错误码，并返回匹配的 HTTP 状态码；
// - 非业务错误一律按 500 处理，避免泄露内部细节。
func Write(w http.ResponseWriter, r *http.Request, data any, err error) {
	resp := OkData(data)
	status := http.StatusOK
	if err != nil {
		resp = Fail(err)
		status = xerr.HttpStatus(err)
	}
	resp.RequestId = trace.TraceIDFromContext(r.Context())
	writeJSON(w, status, resp)
}

func writeJSON(w http.ResponseWriter, status int, resp *R) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = jsonEncode(w, resp)
}
