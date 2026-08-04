package result

import (
	"encoding/json"
	"io"
)

// jsonEncode 输出 JSON。这里保留独立函数，便于后续替换为更高效的序列化实现。
func jsonEncode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
