package handler

import (
	"regexp"
	"strconv"
	"strings"

	"common/xerr"
	"user/internal/types"
)

var cellPhoneRe = regexp.MustCompile(`^1\d{10}$`)

// ValidateLoginForm 校验登录表单
func ValidateLoginForm(req *types.LoginFormDTO) error {
	if req.Type <= 0 {
		return xerr.BadRequestf("用户类型不能为空")
	}
	validTypes := map[int64]struct{}{1: {}, 2: {}, 3: {}}
	if _, ok := validTypes[req.Type]; !ok {
		return xerr.BadRequestf("非法的用户类型")
	}
	account := strings.TrimSpace(req.CellPhone)
	if account == "" {
		account = strings.TrimSpace(req.Username)
	}
	if account == "" {
		return xerr.BadRequestf("手机号或用户名不能为空")
	}
	if req.Password == "" {
		return xerr.BadRequestf("密码不能为空")
	}
	return nil
}

// ValidateStudentForm 校验学员表单
func ValidateStudentForm(req *types.StudentFormDTO) error {
	if req.CellPhone == "" {
		return xerr.BadRequestf("手机号不能为空")
	}
	if !cellPhoneRe.MatchString(req.CellPhone) {
		return xerr.BadRequestf("手机号格式不正确")
	}
	if req.Password == "" {
		return xerr.BadRequestf("密码不能为空")
	}
	return nil
}

// ValidateUserIdReq 校验 UserIdReq
func ValidateUserIdReq(req *types.UserIdReq) error {
	if req.Id <= 0 {
		return xerr.BadRequestf("用户 id 必须大于 0")
	}
	return nil
}

// ValidateUserDTO 校验 UserDTO（非空字段校验）
func ValidateUserDTO(req *types.UserDTO) error {
	if req.Username == "" {
		return xerr.BadRequestf("用户名不能为空")
	}
	if req.CellPhone == "" {
		return xerr.BadRequestf("手机号不能为空")
	}
	if req.Type == 2 && !cellPhoneRe.MatchString(req.CellPhone) {
		return xerr.BadRequestf("手机号格式不正确")
	}
	return nil
}

// ParsePathID 从路径参数解析 :id
func ParsePathID(pathIDStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(pathIDStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, xerr.New(xerr.CodeBadRequest, "路径参数 id 非法")
	}
	return id, nil
}