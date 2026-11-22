package courselogic

import (
	"database/sql"

	"tjxt/pkg/timex"
	"tjxt/rpc/course/internal/model"
	"tjxt/rpc/course/pb"
)

// formatNullInt64 取 sql.NullInt64 有效值，无效返回 0。
func formatNullInt64(v sql.NullInt64) int64 {
	if v.Valid {
		return v.Int64
	}
	return 0
}

// buildCategoryTree 由扁平分类列表构造树形结构（多级）。
func buildCategoryTree(list []*model.Category) []*pb.CategoryNode {
	nodes := make(map[int64]*pb.CategoryNode, len(list))
	roots := make([]*pb.CategoryNode, 0, len(list))
	for _, c := range list {
		nodes[c.Id] = &pb.CategoryNode{
			Id:       c.Id,
			Name:     c.Name,
			ParentId: c.ParentId,
			Level:    int32(c.Level),
			Index:    int32(c.Priority),
			Status:   int32(c.Status),
		}
	}
	for _, n := range nodes {
		if n.ParentId != 0 {
			if p, ok := nodes[n.ParentId]; ok {
				p.Children = append(p.Children, n)
				continue
			}
		}
		roots = append(roots, n)
	}
	return roots
}

// toCategoryInfo 分类实体转 pb（courseNum / thirdCategoryNum 由调用方查询后传入）
func toCategoryInfo(c *model.Category, courseNum, thirdNum int64) *pb.CategoryInfo {
	return &pb.CategoryInfo{
		Id:               c.Id,
		Name:             c.Name,
		ParentId:         c.ParentId,
		Level:            int32(c.Level),
		Priority:         int32(c.Priority),
		Status:           int32(c.Status),
		StatusDesc:       categoryStatusDesc(int64(c.Status)),
		CourseNum:        courseNum,
		ThirdCategoryNum: thirdNum,
		CreateTime:       timex.Format(c.CreateTime),
		UpdateTime:       timex.Format(c.UpdateTime),
	}
}

// categoryNameMap 构造 id -> 分类 映射，便于回溯一/二级分类名称。
func categoryNameMap(list []*model.Category) map[int64]*model.Category {
	m := make(map[int64]*model.Category, len(list))
	for _, c := range list {
		m[c.Id] = c
	}
	return m
}
