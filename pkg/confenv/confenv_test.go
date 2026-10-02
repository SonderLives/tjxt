package confenv

import (
	"os"
	"testing"
)

func TestExpandWithDefault(t *testing.T) {
	os.Setenv("TJXT_TEST_SECRET", "prod-secret")
	defer os.Unsetenv("TJXT_TEST_SECRET")

	got := string(Expand([]byte(`secret: "${TJXT_TEST_SECRET|fallback}"`)))
	if got != `secret: "prod-secret"` {
		t.Fatalf("env set: %q", got)
	}

	got = string(Expand([]byte(`secret: "${TJXT_TEST_UNSET|fallback}"`)))
	if got != `secret: "fallback"` {
		t.Fatalf("env unset: %q", got)
	}
}

func TestExpandEmptyEnvFallsBack(t *testing.T) {
	os.Setenv("TJXT_TEST_EMPTY", "")
	defer os.Unsetenv("TJXT_TEST_EMPTY")

	got := string(Expand([]byte(`db: "root:${TJXT_TEST_EMPTY|0000}@tcp"`)))
	if got != `db: "root:0000@tcp"` {
		t.Fatalf("empty env should fall back: %q", got)
	}
}

func TestExpandPlainContentUntouched(t *testing.T) {
	in := []byte(`name: pay.rpc
DataSource: root:0000@tcp(127.0.0.1:3306)/tj_pay`)
	if got := string(Expand(in)); string(in) != got {
		t.Fatalf("plain content changed: %q", got)
	}
}
