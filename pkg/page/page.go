// Package page 提供分页参数的归一化工具。
package page

// Req 前端通用的分页请求参数，各服务可直接内嵌。
// form tag 供 go-zero httpx.Parse 自动绑定。
type Req struct {
	PageNo   int64  `form:"pageNo,omitempty"`
	PageSize int64  `form:"pageSize,omitempty"`
	IsAsc    bool   `form:"isAsc,omitempty"`
	SortBy   string `form:"sortBy,omitempty"`
}

// Normalize 将分页参数归一化，返回 offset 与 limit。
// 非法值使用默认值：pageNo>=1，1<=pageSize<=100。
func (p *Req) Normalize() (offset, limit int64) {
	return Normalize(p.PageNo, p.PageSize)
}

// Normalize 归一化分页参数。
func Normalize(pageNo, pageSize int64) (offset, limit int64) {
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return (pageNo - 1) * pageSize, pageSize
}

// CalcPages 计算总页数。
func CalcPages(total, pageSize int64) int64 {
	if pageSize <= 0 || total <= 0 {
		return 0
	}
	pages := total / pageSize
	if total%pageSize != 0 {
		pages++
	}
	return pages
}
