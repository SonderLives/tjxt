package authlogic

import (
	"database/sql"
	"sort"
	"time"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/pb"
)

// timeLayout 对外输出的时间格式，与其他微服务保持一致。
const timeLayout = "2006-01-02 15:04:05"

// fixedRoleType 固定角色类型：不允许删除或修改 code。
const fixedRoleType int64 = 0

// strVal 从 sql.NullString 取值，无效时返回空串。
func strVal(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

// int64Val 从 sql.NullInt64 取值，无效时返回 0。
func int64Val(v sql.NullInt64) int64 {
	if v.Valid {
		return v.Int64
	}
	return 0
}

// timeStr 格式化时间，零值返回空串。
func timeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

// nullTimeStr 格式化可空时间，无效或零值返回空串。
func nullTimeStr(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return timeStr(t.Time)
}

// boolToInt 将布尔转为数据库中的 tinyint 表示。
func boolToInt(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

// intToBool 将数据库 tinyint 转为布尔。
func intToBool(v int64) bool {
	return v != 0
}

// toRoleVO model.Role -> pb.RoleVO。
func toRoleVO(r *model.Role) *pb.RoleVO {
	if r == nil {
		return nil
	}
	return &pb.RoleVO{
		Id:         r.Id,
		Code:       r.Code,
		Name:       r.Name,
		Type:       int32(r.Type),
		CreateTime: timeStr(r.CreateTime),
	}
}

// toPrivilegeVO model.Privilege -> pb.PrivilegeVO。
func toPrivilegeVO(p *model.Privilege) *pb.PrivilegeVO {
	if p == nil {
		return nil
	}
	return &pb.PrivilegeVO{
		Id:       p.Id,
		MenuId:   int64Val(p.MenuId),
		Intro:    p.Intro,
		Method:   p.Method,
		Uri:      p.Uri,
		Internal: intToBool(p.Internal),
	}
}

// toLoginRecordVO model.LoginRecord -> pb.LoginRecordVO。
func toLoginRecordVO(r *model.LoginRecord) *pb.LoginRecordVO {
	if r == nil {
		return nil
	}
	return &pb.LoginRecordVO{
		Id:         r.Id,
		UserId:     r.UserId,
		CellPhone:  strVal(r.CellPhone),
		LoginTime:  timeStr(r.LoginTime),
		LogoutTime: nullTimeStr(r.LogoutTime),
		Duration:   r.Duration,
		Ipv4:       r.Ipv4,
	}
}

// toMenuVO model.Menu -> pb.MenuVO（不含 children）。
func toMenuVO(m *model.Menu) *pb.MenuVO {
	if m == nil {
		return nil
	}
	return &pb.MenuVO{
		Id:          m.Id,
		ParentId:    m.ParentId,
		HasChildren: intToBool(m.HasChildren),
		Label:       m.Label,
		Path:        m.Path,
		Icon:        m.Icon,
		Priority:    int32(m.Priority),
		Children:    nil,
	}
}

// buildMenuTree 将扁平菜单列表在内存中组装为多级树。
//
// 一次遍历建索引、二次遍历挂载父子关系，复杂度 O(n)，避免逐层查库。
// 父节点缺失（被删除或数据不一致）的节点会被提升为根节点，防止数据丢失。
func buildMenuTree(menus []*model.Menu) []*pb.MenuVO {
	if len(menus) == 0 {
		return []*pb.MenuVO{}
	}

	nodes := make(map[int64]*pb.MenuVO, len(menus))
	order := make([]int64, 0, len(menus))
	for _, m := range menus {
		nodes[m.Id] = toMenuVO(m)
		order = append(order, m.Id)
	}

	roots := make([]*pb.MenuVO, 0, len(menus))
	for _, id := range order {
		node := nodes[id]
		parent, ok := nodes[node.ParentId]
		if node.ParentId <= 0 || !ok {
			roots = append(roots, node)
			continue
		}
		parent.Children = append(parent.Children, node)
	}

	// has_children 以实际挂载结果为准，纠正历史脏数据。
	for _, node := range nodes {
		node.HasChildren = len(node.Children) > 0
		sortMenus(node.Children)
	}
	sortMenus(roots)
	return roots
}

// sortMenus 按 priority 升序、id 升序稳定排序，保证树的输出顺序确定。
func sortMenus(list []*pb.MenuVO) {
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].Priority != list[j].Priority {
			return list[i].Priority < list[j].Priority
		}
		return list[i].Id < list[j].Id
	})
}
