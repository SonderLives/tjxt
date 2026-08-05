// Package convert 负责 user API 层与 RPC 层之间的协议转换（pb <-> types）。
// 业务手写辅助代码，goctl 不会覆盖此包。
package convert

import (
	"tjxt/apps/user/api/internal/types"
	pb "tjxt/apps/user/rpc/pb"
)

// ===== RPC -> API 视图对象 =====

func ToIdVO(in *pb.IdResponse) *types.IdVO {
	if in == nil {
		return nil
	}
	return &types.IdVO{Id: in.Id}
}

func ToBoolVO(in *pb.BoolResponse) *types.BoolVO {
	if in == nil {
		return nil
	}
	return &types.BoolVO{Result: in.Result}
}

func ToUserDTO(in *pb.UserDTO) *types.UserDTO {
	if in == nil {
		return nil
	}
	return &types.UserDTO{
		Id:        in.Id,
		Username:  in.Username,
		CellPhone: in.CellPhone,
		Type:      in.Type,
		Name:      in.Name,
		Gender:    in.Gender,
		Icon:      in.Icon,
		Email:     in.Email,
		Qq:        in.Qq,
		Job:       in.Job,
		Province:  in.Province,
		City:      in.City,
		District:  in.District,
		Intro:     in.Intro,
		Photo:     in.Photo,
		RoleId:    in.RoleId,
	}
}

func ToUserDetailVO(in *pb.UserDetailVO) *types.UserDetailVO {
	if in == nil {
		return nil
	}
	return &types.UserDetailVO{
		Id:         in.Id,
		Username:   in.Username,
		CellPhone:  in.CellPhone,
		Name:       in.Name,
		Gender:     in.Gender,
		Icon:       in.Icon,
		Email:      in.Email,
		Qq:         in.Qq,
		Province:   in.Province,
		City:       in.City,
		District:   in.District,
		Intro:      in.Intro,
		RoleName:   in.RoleName,
		CreateTime: in.CreateTime,
	}
}

func toStudentPageVO(in *pb.StudentPageVO) *types.StudentPageVO {
	if in == nil {
		return nil
	}
	return &types.StudentPageVO{
		Id:           in.Id,
		Name:         in.Name,
		CellPhone:    in.CellPhone,
		Icon:         in.Icon,
		Gender:       in.Gender,
		Status:       in.Status,
		CourseAmount: in.CourseAmount,
		CreateTime:   in.CreateTime,
	}
}

func ToStudentPageVOList(out *pb.StudentPageResponse) *types.StudentPageVOList {
	if out == nil {
		return nil
	}
	list := make([]types.StudentPageVO, 0, len(out.List))
	for _, v := range out.List {
		list = append(list, *toStudentPageVO(v))
	}
	return &types.StudentPageVOList{Total: out.Total, Pages: out.Pages, List: list}
}

func toTeacherPageVO(in *pb.TeacherPageVO) *types.TeacherPageVO {
	if in == nil {
		return nil
	}
	return &types.TeacherPageVO{
		Id:                 in.Id,
		Name:               in.Name,
		CellPhone:          in.CellPhone,
		Icon:               in.Icon,
		Photo:              in.Photo,
		Job:                in.Job,
		Intro:              in.Intro,
		Status:             in.Status,
		CourseAmount:       in.CourseAmount,
		ExamQuestionAmount: in.ExamQuestionAmount,
		CreateTime:         in.CreateTime,
	}
}

func ToTeacherPageVOList(out *pb.TeacherPageResponse) *types.TeacherPageVOList {
	if out == nil {
		return nil
	}
	list := make([]types.TeacherPageVO, 0, len(out.List))
	for _, v := range out.List {
		list = append(list, *toTeacherPageVO(v))
	}
	return &types.TeacherPageVOList{Total: out.Total, Pages: out.Pages, List: list}
}

func toStaffVO(in *pb.StaffVO) *types.StaffVO {
	if in == nil {
		return nil
	}
	return &types.StaffVO{
		Id:         in.Id,
		Name:       in.Name,
		CellPhone:  in.CellPhone,
		Icon:       in.Icon,
		Status:     in.Status,
		RoleId:     in.RoleId,
		RoleName:   in.RoleName,
		CreateTime: in.CreateTime,
	}
}

func ToStaffVOList(out *pb.StaffPageResponse) *types.StaffVOList {
	if out == nil {
		return nil
	}
	list := make([]types.StaffVO, 0, len(out.List))
	for _, v := range out.List {
		list = append(list, *toStaffVO(v))
	}
	return &types.StaffVOList{Total: out.Total, Pages: out.Pages, List: list}
}

// ===== API -> RPC 请求 =====

func FromStudentFormReq(in *types.StudentFormReq) *pb.StudentFormRequest {
	return &pb.StudentFormRequest{
		CellPhone: in.CellPhone,
		Code:      in.Code,
		Password:  in.Password,
	}
}

func FromUserDTO(in *types.UserDTO) *pb.UserDTO {
	return &pb.UserDTO{
		Id:        in.Id,
		Username:  in.Username,
		CellPhone: in.CellPhone,
		Type:      in.Type,
		Name:      in.Name,
		Gender:    in.Gender,
		Icon:      in.Icon,
		Email:     in.Email,
		Qq:        in.Qq,
		Job:       in.Job,
		Province:  in.Province,
		City:      in.City,
		District:  in.District,
		Intro:     in.Intro,
		Photo:     in.Photo,
		RoleId:    in.RoleId,
	}
}

func FromUserFormReq(in *types.UserFormReq) *pb.UserFormRequest {
	return &pb.UserFormRequest{
		Username:   in.Username,
		CellPhone:  in.CellPhone,
		Type:       in.Type,
		Name:       in.Name,
		Gender:     in.Gender,
		Icon:       in.Icon,
		Email:      in.Email,
		Qq:         in.Qq,
		Job:        in.Job,
		Province:   in.Province,
		City:       in.City,
		District:   in.District,
		Intro:      in.Intro,
		Photo:      in.Photo,
		RoleId:     in.RoleId,
		OldPassword: in.OldPassword,
		Password:   in.Password,
	}
}

func FromUserUpdateReq(in *types.UserUpdateReq) *pb.UserDTO {
	return &pb.UserDTO{
		Id:        in.Id,
		Username:  in.Username,
		CellPhone: in.CellPhone,
		Type:      in.Type,
		Name:      in.Name,
		Gender:    in.Gender,
		Icon:      in.Icon,
		Email:     in.Email,
		Qq:        in.Qq,
		Job:       in.Job,
		Province:  in.Province,
		City:      in.City,
		District:  in.District,
		Intro:     in.Intro,
		Photo:     in.Photo,
		RoleId:    in.RoleId,
	}
}

func FromCheckCellPhoneReq(in *types.CheckCellPhoneReq) *pb.CheckCellPhoneRequest {
	return &pb.CheckCellPhoneRequest{CellPhone: in.CellPhone}
}

func FromUpdateStatusReq(in *types.UpdateStatusReq, operator int64) *pb.UpdateStatusRequest {
	return &pb.UpdateStatusRequest{
		UserId:  in.Id,
		Status:  in.Status,
		Operator: operator,
	}
}

func FromUserPageReq(in *types.UserPageReq) *pb.UserPageRequest {
	return &pb.UserPageRequest{
		PageNo:   in.PageNo,
		PageSize: in.PageSize,
		SortBy:   in.SortBy,
		IsAsc:    in.IsAsc,
		Name:     in.Name,
		Phone:    in.Phone,
		Status:   in.Status,
	}
}
