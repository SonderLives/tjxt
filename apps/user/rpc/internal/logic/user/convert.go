package userlogic

import (
	"database/sql"

	"tjxt/apps/user/rpc/internal/model"
	"tjxt/apps/user/rpc/pb"
)

// nullStr 将字符串转为 sql.NullString，空字符串视为 NULL。
func nullStr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// applyStr 合并更新：next 为空则保留原值，否则以新值覆盖。
func applyStr(cur sql.NullString, next string) sql.NullString {
	if next == "" {
		return cur
	}
	return nullStr(next)
}

// strVal 从 sql.NullString 取出字符串，无效时返回空。
func strVal(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

// defaultInitialPassword 管理端新增用户 / 重置密码时的默认口令，正式环境应改为短信或随机下发。
const defaultInitialPassword = "123456"

// toUserDTOFromView 联合视图 -> UserDTO（GetUsersByIds 复用，避免 N+1 查询）。
func toUserDTOFromView(v *model.UserWithDetail) *pb.UserDTO {
	return &pb.UserDTO{
		Id:        v.Id,
		CellPhone: v.CellPhone,
		Type:      int32(v.Type),
		Name:      strVal(v.Name),
		Gender:    int32(v.Gender),
		Icon:      strVal(v.Icon),
		Email:     strVal(v.Email),
		Qq:        strVal(v.Qq),
		Job:       strVal(v.Job),
		Province:  strVal(v.Province),
		City:      strVal(v.City),
		District:  strVal(v.District),
		Intro:     strVal(v.Intro),
		Photo:     strVal(v.Photo),
		RoleId:    v.RoleId,
	}
}

// toUserDTO user + detail -> UserDTO.
func toUserDTO(u *model.User, d *model.UserDetail) *pb.UserDTO {
	out := &pb.UserDTO{
		Id:        u.Id,
		Username:  u.Username,
		CellPhone: u.CellPhone,
		Type:      int32(u.Type),
	}
	if d != nil {
		out.Name = d.Name
		out.Gender = int32(d.Gender)
		out.Icon = strVal(d.Icon)
		out.Email = strVal(d.Email)
		out.Qq = strVal(d.Qq)
		out.Job = strVal(d.Job)
		out.Province = strVal(d.Province)
		out.City = strVal(d.City)
		out.District = strVal(d.District)
		out.Intro = strVal(d.Intro)
		out.Photo = strVal(d.Photo)
		out.RoleId = d.RoleId
	}
	return out
}

// toUserDetailVO user + detail -> UserDetailVO。
// role_name 归属 auth 服务的 role 表，跨服务查询留待后续，暂留空。
func toUserDetailVO(u *model.User, d *model.UserDetail) *pb.UserDetailVO {
	out := &pb.UserDetailVO{
		Id:        u.Id,
		Username:  u.Username,
		CellPhone: u.CellPhone,
	}
	if d != nil {
		out.Name = d.Name
		out.Gender = int32(d.Gender)
		out.Icon = strVal(d.Icon)
		out.Email = strVal(d.Email)
		out.Qq = strVal(d.Qq)
		out.Province = strVal(d.Province)
		out.City = strVal(d.City)
		out.District = strVal(d.District)
		out.Intro = strVal(d.Intro)
	}
	out.CreateTime = u.CreateTime.Format("2006-01-02 15:04:05")
	return out
}

// toStudentPageVO 联合视图 -> StudentPageVO。
func toStudentPageVO(v *model.UserWithDetail) *pb.StudentPageVO {
	return &pb.StudentPageVO{
		Id:           v.Id,
		Name:         strVal(v.Name),
		CellPhone:    v.CellPhone,
		Icon:         strVal(v.Icon),
		Gender:       int32(v.Gender),
		Status:       int32(v.Status),
		CourseAmount: v.CourseAmount,
		CreateTime:   v.CreateTime.Format("2006-01-02 15:04:05"),
	}
}

// toTeacherPageVO 联合视图 -> TeacherPageVO。
// exam_question_amount 归属 exam 服务；course_amount 此处取 user_detail 的购课数（权威值由 course 服务计算），
// 跨服务聚合留待后续。
func toTeacherPageVO(v *model.UserWithDetail) *pb.TeacherPageVO {
	return &pb.TeacherPageVO{
		Id:           v.Id,
		Name:         strVal(v.Name),
		CellPhone:    v.CellPhone,
		Icon:         strVal(v.Icon),
		Photo:        strVal(v.Photo),
		Job:          strVal(v.Job),
		Intro:        strVal(v.Intro),
		Status:       int32(v.Status),
		CourseAmount: v.CourseAmount,
		CreateTime:   v.CreateTime.Format("2006-01-02 15:04:05"),
	}
}

// toStaffVO 联合视图 -> StaffVO。
// role_name 归属 auth 服务的 role 表，跨服务查询留待后续，暂留空。
func toStaffVO(v *model.UserWithDetail) *pb.StaffVO {
	return &pb.StaffVO{
		Id:         v.Id,
		Name:       strVal(v.Name),
		CellPhone:  v.CellPhone,
		Icon:       strVal(v.Icon),
		Status:     int32(v.Status),
		RoleId:     v.RoleId,
		CreateTime: v.CreateTime.Format("2006-01-02 15:04:05"),
	}
}

// sortClause 将用户可控的排序字段映射为白名单列，杜绝 SQL 注入。
// 空 sortBy 默认按注册时间倒序。
func sortClause(sortBy string, isAsc bool) (string, string) {
	var col string
	switch sortBy {
	case "name":
		col = "d.`name`"
	case "status":
		col = "u.`status`"
	case "id":
		col = "u.`id`"
	default:
		col = "u.`create_time`"
	}
	dir := "DESC"
	if isAsc {
		dir = "ASC"
	}
	return col, dir
}
