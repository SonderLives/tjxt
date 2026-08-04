package types

// ============ 认证 ============

// LoginFormDTO 登录请求
type LoginFormDTO struct {
	CellPhone  string `json:"cellPhone,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password"`
	Type       int64  `json:"type"`
	RememberMe bool   `json:"rememberMe,omitempty"`
}

// ============ 学员 ============

// StudentFormDTO 学员注册/改密请求
type StudentFormDTO struct {
	CellPhone string `json:"cellPhone"`
	Code      string `json:"code,omitempty"`
	Password  string `json:"password"`
}

// ============ 分页请求 ============

// UserPageReq 分页查询请求（学员/教师/员工通用）
type UserPageReq struct {
	PageNo   int64  `form:"pageNo,omitempty"`
	PageSize int64  `form:"pageSize,omitempty"`
	IsAsc    bool   `form:"isAsc,omitempty"`
	SortBy   string `form:"sortBy,omitempty"`
	Name     string `form:"name,omitempty"`
	Phone    string `form:"phone,omitempty"`
	Status   int64  `form:"status,omitempty"`
}

// UserCheckCellphoneReq 校验手机号请求
type UserCheckCellphoneReq struct {
	Cellphone string `form:"cellphone"`
}

// UserIdReq 路径用户 id 请求
type UserIdReq struct {
	Id int64 `path:"id"`
}

// UserStatusReq 修改用户状态请求
type UserStatusReq struct {
	Id     int64 `path:"id"`
	Status int64 `path:"status"`
}

// ============ 用户 DTO ============

// UserDTO 用户数据（创建/更新/查询返回）
type UserDTO struct {
	Id        int64  `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	CellPhone string `json:"cellPhone,omitempty"`
	Name      string `json:"name,omitempty"`
	Gender    int64  `json:"gender,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Email     string `json:"email,omitempty"`
	Qq        string `json:"qq,omitempty"`
	Job       string `json:"job,omitempty"`
	Province  string `json:"province,omitempty"`
	City      string `json:"city,omitempty"`
	District  string `json:"district,omitempty"`
	Intro     string `json:"intro,omitempty"`
	Photo     string `json:"photo,omitempty"`
	RoleId    int64  `json:"roleId,omitempty"`
	Type      int64  `json:"type,omitempty"`
}

// UserFormDTO 当前用户更新请求
type UserFormDTO struct {
	UserDTO
	OldPassword string `json:"oldPassword,omitempty"`
	Password    string `json:"password,omitempty"`
}

// ============ 用户 VO ============

// UserDetailVO 当前用户信息
type UserDetailVO struct {
	Id         int64  `json:"id"`
	Username   string `json:"username"`
	CellPhone  string `json:"cellPhone"`
	Name       string `json:"name"`
	Gender     int64  `json:"gender"`
	Icon       string `json:"icon"`
	Email      string `json:"email"`
	Qq         string `json:"qq"`
	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	Intro      string `json:"intro"`
	RoleName   string `json:"roleName"`
	CreateTime string `json:"createTime"`
}

// StudentPageVO 学员分页条目
type StudentPageVO struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	CellPhone    string `json:"cellPhone"`
	Gender       int64  `json:"gender"`
	Icon         string `json:"icon"`
	Status       int64  `json:"status"`
	CourseAmount int64  `json:"courseAmount"`
	CreateTime   string `json:"createTime"`
}

// TeacherPageVO 教师分页条目
type TeacherPageVO struct {
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	CellPhone          string `json:"cellPhone"`
	Icon               string `json:"icon"`
	Photo              string `json:"photo"`
	Job                string `json:"job"`
	Intro              string `json:"intro"`
	Status             int64  `json:"status"`
	CourseAmount       int64  `json:"courseAmount"`
	ExamQuestionAmount int64  `json:"examQuestionAmount"`
	CreateTime         string `json:"createTime"`
}

// StaffVO 员工分页条目
type StaffVO struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	CellPhone  string `json:"cellPhone"`
	Icon       string `json:"icon"`
	RoleId     int64  `json:"roleId"`
	RoleName   string `json:"roleName"`
	Status     int64  `json:"status"`
	CreateTime string `json:"createTime"`
}

// Page 通用分页响应
type Page struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Pages int64 `json:"pages"`
}
