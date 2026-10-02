package page

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		pageNo, pageSize, wantOffset, wantLimit int64
	}{
		{1, 10, 0, 10},   // 正常
		{0, 10, 0, 10},   // 页码归一
		{-3, 0, 0, 10},   // 默认页大小
		{2, 10, 10, 10},  // 第二页
		{3, 100, 200, 100},
		{1, 500, 0, 100}, // 上限截断
	}
	for _, c := range cases {
		offset, limit := Normalize(c.pageNo, c.pageSize)
		if offset != c.wantOffset || limit != c.wantLimit {
			t.Fatalf("Normalize(%d,%d) = (%d,%d), want (%d,%d)",
				c.pageNo, c.pageSize, offset, limit, c.wantOffset, c.wantLimit)
		}
	}
}

func TestReqNormalize(t *testing.T) {
	r := Req{PageNo: 4, PageSize: 25}
	offset, limit := r.Normalize()
	if offset != 75 || limit != 25 {
		t.Fatalf("Req.Normalize = (%d,%d), want (75,25)", offset, limit)
	}
}

func TestCalcPages(t *testing.T) {
	cases := []struct {
		total, pageSize, want int64
	}{
		{0, 10, 0},
		{10, 0, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{101, 20, 6},
	}
	for _, c := range cases {
		if got := CalcPages(c.total, c.pageSize); got != c.want {
			t.Fatalf("CalcPages(%d,%d) = %d, want %d", c.total, c.pageSize, got, c.want)
		}
	}
}
