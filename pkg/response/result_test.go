package response_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tjxt/pkg/response"
	"tjxt/pkg/xerr"
)

func TestWriteSuccess(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	response.Write(w, r, map[string]any{"k": "v"}, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	body := w.Body.String()
	// requestId 在无 trace 的测试环境为空串，但字段必须存在
	for _, want := range []string{`"code":200`, `"msg":"OK"`, `"requestId":`, `"data":{"k":"v"}`} {
		if !strings.Contains(body, want) {
			t.Fatalf("body %s missing %s", body, want)
		}
	}
}

func TestWriteBizError(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()

	response.Write(w, r, nil, xerr.New(xerr.CodeNotFound, "课程不存在"))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `"code":404`) || !strings.Contains(body, "课程不存在") {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWritePlainErrorMaskedAs500(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()

	response.Write(w, r, nil, errors.New("internal secret detail"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `"code":500`) {
		t.Fatalf("unexpected body: %s", body)
	}
	// 非业务错误的内部细节不得透出
	if strings.Contains(body, "secret") {
		t.Fatalf("internal detail leaked: %s", body)
	}
}

func TestFailNilErrReturnsOK(t *testing.T) {
	if got := response.Fail(nil); got.Code != 200 {
		t.Fatalf("Fail(nil).Code = %d", got.Code)
	}
}

func TestPageStruct(t *testing.T) {
	p := response.Page{List: []int{1, 2}, Total: 2, Pages: 1}
	if p.Total != 2 || p.Pages != 1 {
		t.Fatalf("unexpected page: %+v", p)
	}
}
