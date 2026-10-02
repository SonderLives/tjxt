// Package confenv 提供支持环境变量覆盖的配置加载。
//
// 语法：yaml 值中写 ${VAR} 或 ${VAR|default}。
//   - 环境变量存在且非空 → 用环境变量；
//   - 否则 → 用 default（未提供 default 则展开为空串）。
//
// 生产通过环境变量注入敏感配置（DB 密码/JWT Secret/MQ 密码），
// 本地开发依赖 yaml 里的 default，无需额外设置。
package confenv

import (
	"os"
	"regexp"

	"github.com/zeromicro/go-zero/core/conf"
)

var envRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(?:\|([^}]*))?\}`)

// Expand 展开 content 中的 ${VAR} / ${VAR|default}。
func Expand(content []byte) []byte {
	return envRe.ReplaceAllFunc(content, func(m []byte) []byte {
		sub := envRe.FindSubmatch(m)
		if v, ok := os.LookupEnv(string(sub[1])); ok && v != "" {
			return []byte(v)
		}
		if len(sub) > 2 {
			return sub[2]
		}
		return nil
	})
}

// MustLoad 读取 yaml 配置，先做环境变量展开再解析；失败时 panic。
func MustLoad(path string, v any) {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := conf.LoadFromYamlBytes(Expand(content), v); err != nil {
		panic(err)
	}
}
