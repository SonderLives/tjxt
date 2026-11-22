package promotionlogic

import (
	"crypto/rand"
	"math/big"
	"strings"

	"tjxt/pkg/errorx"
)

// codeAlphabet 兑换码字符集，剔除易混淆的 0/O/1/I。
const codeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

// codeLength 兑换码长度。
const codeLength = 12

// generateCodes 生成 n 个互不相同的兑换码；DB 上 uk_code 唯一索引兜底防重。
func generateCodes(n int64) ([]string, error) {
	if n <= 0 {
		return nil, nil
	}
	seen := make(map[string]struct{}, n)
	codes := make([]string, 0, n)
	max := big.NewInt(int64(len(codeAlphabet)))

	for int64(len(codes)) < n {
		var sb strings.Builder
		sb.Grow(codeLength)
		for i := 0; i < codeLength; i++ {
			idx, err := rand.Int(rand.Reader, max)
			if err != nil {
				return nil, errorx.Wrap(err, errorx.CodeInternal, "生成兑换码失败")
			}
			sb.WriteByte(codeAlphabet[idx.Int64()])
		}
		code := sb.String()
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		codes = append(codes, code)
	}
	return codes, nil
}
