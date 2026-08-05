package logic

import (
	"database/sql"
	"time"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/utils/idgen"
)

// 分类状态
const (
	CategoryStatusNormal  = 1
	CategoryStatusDisable = 2
)

// 课程状态（与 course / course_draft 表注释一致）
const (
	CourseStatusPending   = 1 // 待上架
	CourseStatusUpShelf   = 2 // 已上架
	CourseStatusDownShelf = 3 // 下架
	CourseStatusFinished  = 4 // 已完结
)

// 目录类型
const (
	CatalogueTypeChapter  = 1 // 章
	CatalogueTypeSection  = 2 // 节
	CatalogueTypePractice = 3 // 测试/练习
)

// nextID 雪花算法生成分布式 ID（course 等草稿表使用，category 使用自增）
func nextID() int64 { return idgen.NextID() }

// isNotFound 判断是否是数据不存在错误
func isNotFound(err error) bool {
	return err == sql.ErrNoRows || err == model.ErrNotFound
}

// courseStatusDesc 课程状态中文描述
func courseStatusDesc(status int64) string {
	switch status {
	case CourseStatusPending:
		return "待上架"
	case CourseStatusUpShelf:
		return "已上架"
	case CourseStatusDownShelf:
		return "已下架"
	case CourseStatusFinished:
		return "已完结"
	default:
		return ""
	}
}

// categoryStatusDesc 分类状态中文描述
func categoryStatusDesc(status int64) string {
	switch status {
	case CategoryStatusNormal:
		return "正常"
	case CategoryStatusDisable:
		return "禁用"
	default:
		return ""
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format("2006-01-02 15:04:05")
}

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
		CreateTime:       formatTime(c.CreateTime),
		UpdateTime:       formatTime(c.UpdateTime),
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
