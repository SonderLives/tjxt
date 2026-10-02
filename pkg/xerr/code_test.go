package xerr

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestNewAndErrorString(t *testing.T) {
	e := New(CodeNotFound, "课程不存在")
	if e.Code != CodeNotFound || e.Msg != "课程不存在" {
		t.Fatalf("unexpected error fields: %+v", e)
	}
	if got, want := e.Error(), "code=404 msg=课程不存在"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestWrapKeepsCause(t *testing.T) {
	cause := errors.New("db timeout")
	e := Wrap(cause, CodeInternal, "查询订单失败")
	if !errors.Is(e, cause) {
		t.Fatal("errors.Is should find wrapped cause")
	}
	if e.Error() != "code=500 msg=查询订单失败 cause=db timeout" {
		t.Fatalf("Error() = %q", e.Error())
	}
	if Wrap(nil, CodeInternal, "x") != nil {
		t.Fatal("Wrap(nil) should return nil")
	}
}

func TestCodeOfFollowsChain(t *testing.T) {
	biz := New(CodeForbidden, "无权操作")
	wrapped := fmt.Errorf("handler: %w", fmt.Errorf("logic: %w", biz))
	if got := CodeOf(wrapped); got != CodeForbidden {
		t.Fatalf("CodeOf = %d, want %d", got, CodeForbidden)
	}
	if got := MsgOf(wrapped); got != "无权操作" {
		t.Fatalf("MsgOf = %q", got)
	}
}

func TestCodeOfPlainErrorIsInternal(t *testing.T) {
	if got := CodeOf(errors.New("boom")); got != CodeInternal {
		t.Fatalf("CodeOf = %d, want %d", got, CodeInternal)
	}
	if got := MsgOf(errors.New("boom")); got != MsgInternal {
		t.Fatalf("MsgOf = %q, want %q", got, MsgInternal)
	}
}

func TestHttpStatusMapping(t *testing.T) {
	cases := map[Code]int{
		CodeSuccess:            http.StatusOK,
		CodeBadRequest:         http.StatusBadRequest,
		CodeUnauthorized:       http.StatusUnauthorized,
		CodeForbidden:          http.StatusForbidden,
		CodeNotFound:           http.StatusNotFound,
		CodeConflict:           http.StatusConflict,
		CodeTooManyRequests:    http.StatusTooManyRequests,
		CodeServiceUnavailable: http.StatusServiceUnavailable,
		Code(999):              http.StatusInternalServerError,
	}
	for code, want := range cases {
		if got := HttpStatus(New(code, "x")); got != want {
			t.Fatalf("HttpStatus(code=%d) = %d, want %d", code, got, want)
		}
	}
}

func TestHelperConstructors(t *testing.T) {
	if e := BadRequestf("字段 %s 非法", "id"); e.Code != CodeBadRequest {
		t.Fatalf("BadRequestf code = %d", e.Code)
	}
	if e := Unauthorized(""); e.Msg != MsgUnauthorized {
		t.Fatalf("Unauthorized default msg = %q", e.Msg)
	}
	if e := NotFound(""); e.Code != CodeNotFound || e.Msg != MsgNotFound {
		t.Fatalf("NotFound default = %+v", e)
	}
}
