package types

// ============ 认证 ============

// LoginFormDTO 登录请求
type LoginFormDTO struct {
	CellPhone  string `json:"cellPhone,optional"`
	Username   string `json:"username,optional"`
	Password   string `json:"password"`
	Type       int64  `json:"type"`
	RememberMe bool   `json:"rememberMe,optional"`
}

// ============ 学员 ============

// StudentFormDTO 学员注册/改密请求
type StudentFormDTO struct {
	CellPhone string `json:"cellPhone"`
	Code      string `json:"code,optional"`
	Password  string `json:"password"`
}

// ============ 分页请求 ============

// UserPageReq 分页查询请求（学员/教师/员工通用）
type UserPageReq struct {
	PageNo   int64  `form:"pageNo,optional"`
	PageSize int64  `form:"pageSize,optional"`
	IsAsc    bool   `form:"isAsc,optional"`
	SortBy   string `form:"sortBy,optional"`
	Name     string `form:"name,optional"`
	Phone    string `form:"phone,optional"`
	Status   int64  `form:"status,optional"`
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
	Id        int64  `json:"id,optional"`
	Username  string `json:"username,optional"`
	CellPhone string `json:"cellPhone,optional"`
	Name      string `json:"name,optional"`
	Gender    int64  `json:"gender,optional"`
	Icon      string `json:"icon,optional"`
	Email     string `json:"email,optional"`
	Qq        string `json:"qq,optional"`
	Job       string `json:"job,optional"`
	Province  string `json:"province,optional"`
	City      string `json:"city,optional"`
	District  string `json:"district,optional"`
	Intro     string `json:"intro,optional"`
	Photo     string `json:"photo,optional"`
	RoleId    int64  `json:"roleId,optional"`
	Type      int64  `json:"type,optional"`
}

// UserFormDTO 当前用户更新请求
type UserFormDTO struct {
	UserDTO
	OldPassword string `json:"oldPassword,optional"`
	Password    string `json:"password,optional"`
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
