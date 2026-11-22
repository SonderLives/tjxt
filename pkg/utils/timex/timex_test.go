package timex

import (
	"database/sql"
	"testing"
	"time"
)

func TestLayout(t *testing.T) {
	// 布局必须严格是 yyyy-MM-dd HH:mm:ss（README 与服务重构依赖该约定）
	if Layout != "2006-01-02 15:04:05" {
		t.Fatalf("Layout = %q", Layout)
	}
}

func TestFormat(t *testing.T) {
	cases := []struct {
		name string
		in   time.Time
		want string
	}{
		{
			name: "零值返回空串",
			in:   time.Time{},
			want: "",
		},
		{
			name: "常规时间",
			in:   time.Date(2024, 5, 1, 9, 30, 5, 0, time.Local),
			want: "2024-05-01 09:30:05",
		},
		{
			name: "日期与秒数补零",
			in:   time.Date(2024, 1, 2, 3, 4, 5, 0, time.Local),
			want: "2024-01-02 03:04:05",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Format(c.in); got != c.want {
				t.Fatalf("Format = %q, want %q", got, c.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    time.Time
		wantErr bool
	}{
		{
			name: "常规解析（本地时区，壁钟时间不变）",
			in:   "2024-05-01 09:30:05",
			want: time.Date(2024, 5, 1, 9, 30, 5, 0, time.Local),
		},
		{
			name:    "空串报错",
			in:      "",
			wantErr: true,
		},
		{
			name:    "缺时间部分报错",
			in:      "2024-05-01",
			wantErr: true,
		},
		{
			name:    "乱文本报错",
			in:      "not-a-time",
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Parse(c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("Parse(%q) want error, got %v", c.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", c.in, err)
			}
			if !got.Equal(c.want) {
				t.Fatalf("Parse = %v, want %v", got, c.want)
			}
		})
	}
}

func TestFormatParseRoundTrip(t *testing.T) {
	// Format → Parse 往返应无损（同一壁钟时间）
	in := time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local)
	out, err := Parse(Format(in))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Equal(in) {
		t.Fatalf("round trip = %v, want %v", out, in)
	}
}

func TestFormatNullTime(t *testing.T) {
	cases := []struct {
		name string
		in   sql.NullTime
		want string
	}{
		{
			name: "无效值返回空串",
			in:   sql.NullTime{},
			want: "",
		},
		{
			name: "Valid=false 即使带时间也返回空串",
			in:   sql.NullTime{Time: time.Date(2024, 5, 1, 0, 0, 0, 0, time.Local)},
			want: "",
		},
		{
			name: "有效值同 Format",
			in:   sql.NullTime{Time: time.Date(2024, 5, 1, 9, 30, 5, 0, time.Local), Valid: true},
			want: "2024-05-01 09:30:05",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := FormatNullTime(c.in); got != c.want {
				t.Fatalf("FormatNullTime = %q, want %q", got, c.want)
			}
		})
	}
}

func TestParseNullTime(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    sql.NullTime
		wantErr bool
	}{
		{
			name: "空串等价于 NULL",
			in:   "",
			want: sql.NullTime{},
		},
		{
			name: "常规解析",
			in:   "2024-05-01 09:30:05",
			want: sql.NullTime{Time: time.Date(2024, 5, 1, 9, 30, 5, 0, time.Local), Valid: true},
		},
		{
			name:    "非法格式报错",
			in:      "2024-05-01",
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseNullTime(c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("ParseNullTime(%q) want error, got %+v", c.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseNullTime(%q) unexpected error: %v", c.in, err)
			}
			if got.Valid != c.want.Valid || !got.Time.Equal(c.want.Time) {
				t.Fatalf("ParseNullTime = %+v, want %+v", got, c.want)
			}
		})
	}
}
