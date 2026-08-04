// Package types 定义课程微服务对外请求/响应数据结构。
// 字段与 Apifox 契约一致（JSON 命名采用驼峰）。
package types

import "common/page"

// ---------- 课程分类 ----------

// CategoryAddDTO 新增课程分类
type CategoryAddDTO struct {
	Name     string `json:"name"`
	ParentId int64  `json:"parentId"`
	Index    int32  `json:"index"`
}

// CategoryUpdateDTO 更新课程分类
type CategoryUpdateDTO struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Index int32  `json:"index"`
}

// CategoryDisableOrEnableDTO 课程分类启用/禁用
type CategoryDisableOrEnableDTO struct {
	Id     int64 `json:"id"`
	Status int32 `json:"status"`
}

// SimpleCategoryVO 所有课程分类数据（前端分类树）
type SimpleCategoryVO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	Level    int32               `json:"level"`
	ParentId int64               `json:"parentId"`
	Children []*SimpleCategoryVO `json:"children"`
}

// CategoryVO 课程分类信息
type CategoryVO struct {
	Id              int64         `json:"id"`
	Name            string        `json:"name"`
	ParentId        int64         `json:"parentId"`
	Level           int32         `json:"level"`
	Index           int32         `json:"index"`
	Status          int32         `json:"status"`
	StatusDesc      string        `json:"statusDesc"`
	CourseNum       int32         `json:"courseNum"`
	ThirdCategoryNum int32        `json:"thirdCategoryNum"`
	CreateTime      string        `json:"createTime"`
	UpdateTime      string        `json:"updateTime"`
	Children        []*CategoryVO `json:"children"`
}

// CategoryInfoVO 课程分类详情
type CategoryInfoVO struct {
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	Index              int32  `json:"index"`
	CategoryLevel      int32  `json:"categoryLevel"`
	Status             int32  `json:"status"`
	StatusDesc         string `json:"statusDesc"`
	FirstCategoryName  string `json:"firstCategoryName"`
	SecondCategoryName string `json:"secondCategoryName"`
	CreateTime         string `json:"createTime"`
	UpdateTime         string `json:"updateTime"`
}

// ---------- 课程基本信息 ----------

// CourseBaseInfoSaveDTO 课程基本信息保存
type CourseBaseInfoSaveDTO struct {
	Id                int64  `json:"id"`
	Name              string `json:"name"`
	CoverUrl          string `json:"coverUrl"`
	Free              bool   `json:"free"`
	Price             int32  `json:"price"`
	ValidDuration     int32  `json:"validDuration"`
	ThirdCateId       int64  `json:"thirdCateId"`
	Introduce         string `json:"introduce"`
	Detail            string `json:"detail"`
	UsePeople         string `json:"usePeople"`
	PurchaseStartTime string `json:"purchaseStartTime"`
	PurchaseEndTime   string `json:"purchaseEndTime"`
}

// CourseSaveVO 课程保存结果
type CourseSaveVO struct {
	Id int64 `json:"id"`
}

// CourseBaseInfoVO 课程基本信息
type CourseBaseInfoVO struct {
	Id                int64   `json:"id"`
	Name              string  `json:"name"`
	CoverUrl          string  `json:"coverUrl"`
	Price             int32   `json:"price"`
	ValidDuration     int32   `json:"validDuration"`
	Free              bool    `json:"free"`
	FirstCateId       int64   `json:"firstCateId"`
	SecondCateId      int64   `json:"secondCateId"`
	ThirdCateId       int64   `json:"thirdCateId"`
	Status            int32   `json:"status"`
	Step              int32   `json:"step"`
	CanUpdate         bool    `json:"canUpdate"`
	CataTotalNum      int32   `json:"cataTotalNum"`
	CoureScore        float64 `json:"coureScore"`
	CreateTime        string  `json:"createTime"`
	UpdateTime        string  `json:"updateTime"`
	Creater           int64   `json:"creater"`
	Updater           int64   `json:"updater"`
	CreaterName       string  `json:"createrName"`
	UpdaterName       string  `json:"updaterName"`
	CateNames         string  `json:"cateNames"`
	Detail            string  `json:"detail"`
	Introduce         string  `json:"introduce"`
	UsePeople         string  `json:"usePeople"`
	EnrollNum         int32   `json:"enrollNum"`
	StudyNum          int32   `json:"studyNum"`
	RefundNum         int32   `json:"refundNum"`
	RealPayAmount     int32   `json:"realPayAmount"`
	Score             int32   `json:"score"`
	PurchaseStartTime string  `json:"purchaseStartTime"`
	PurchaseEndTime   string  `json:"purchaseEndTime"`
}

// NameExistVO 课程名称是否已存在
type NameExistVO struct {
	Existed bool `json:"existed"`
}

// CourseIdDTO 课程 id
type CourseIdDTO struct {
	Id int64 `json:"id"`
}

// CourseCataIdVO 生成的练习（目录）id
type CourseCataIdVO struct {
	Id int64 `json:"id"`
}

// ---------- 课程目录 ----------

// CataSaveDTO 章节保存模型
type CataSaveDTO struct {
	Id       int64         `json:"id"`
	Index    int32         `json:"index"`
	Name     string        `json:"name"`
	Type     int32         `json:"type"`
	Sections []*CataSaveDTO `json:"sections"`
}

// CataVO 课程目录
type CataVO struct {
	Id                   int64    `json:"id"`
	Index                int32    `json:"index"`
	Name                 string   `json:"name"`
	MediaDuration        int32    `json:"mediaDuration"`
	Trailer              bool     `json:"trailer"`
	MediaName            string   `json:"mediaName"`
	MediaId              int64    `json:"mediaId"`
	Type                 int32    `json:"type"`
	SubjectNum           int32    `json:"subjectNum"`
	TotalScore           int32    `json:"totalScore"`
	CanUpdate            bool     `json:"canUpdate"`
	Sections             []*CataVO `json:"sections"`
	MaxIndexOnShelf      int32    `json:"maxIndexOnShelf"`
	MaxSectionIndexOnShelf int32  `json:"maxSectionIndexOnShelf"`
}

// CataSimpleInfoVO 目录简单信息
type CataSimpleInfoVO struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Index        string `json:"index"`
	CIndex       int32  `json:"cindex"`
	ChapterIndex int32  `json:"chapterIndex"`
}

// CourseMediaSaveDTO 小节媒资保存
type CourseMediaSaveDTO struct {
	CataId        int64  `json:"cataId"`
	MediaId       int64  `json:"mediaId"`
	Trailer       bool   `json:"trailer"`
	VideoName     string `json:"videoName,optional"`
	MediaDuration int32  `json:"mediaDuration,optional"`
}

// ---------- 课程题目 ----------

// CataSubjectDTO 小节/练习与题目关系保存
type CataSubjectDTO struct {
	CataId     int64   `json:"cataId"`
	SubjectIds []int64 `json:"subjectIds"`
}

// SubjectInfo 题目简要信息
type SubjectInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// CataSimpleSubjectVO 小节或练习对应的题目列表
type CataSimpleSubjectVO struct {
	CataId   int64          `json:"cataId"`
	Subjects []*SubjectInfo `json:"subjects"`
}

// ---------- 课程老师 ----------

// TeacherInfo 老师 id 与用户端是否展示
type TeacherInfo struct {
	Id     int64 `json:"id"`
	IsShow bool  `json:"isShow"`
}

// CourseTeacherSaveDTO 保存老师课程关系
type CourseTeacherSaveDTO struct {
	Id       int64          `json:"id"`
	Teachers []*TeacherInfo `json:"teachers"`
}

// CourseTeacherVO 老师课程信息
type CourseTeacherVO struct {
	Id       int64  `json:"id"`
	Icon     string `json:"icon"`
	Photo    string `json:"photo"`
	Name     string `json:"name"`
	Introduce string `json:"introduce"`
	IsShow   bool   `json:"isShow"`
	Job      string `json:"job"`
}

// ---------- 课程详情与进度 ----------

// CourseAndSectionVO 课程和目录及学习进度信息
type CourseAndSectionVO struct {
	Id              int64        `json:"id"`
	Name            string       `json:"name"`
	Sections        int32        `json:"sections"`
	CoverUrl        string       `json:"coverUrl"`
	LessonId        int64        `json:"lessonId"`
	LatestSectionId int64        `json:"latestSectionId"`
	TeacherName     string       `json:"teacherName"`
	TeacherIcon     string       `json:"teacherIcon"`
	Chapters        []*ChapterVO `json:"chapters"`
}

// ChapterVO 章信息
type ChapterVO struct {
	Id            int64        `json:"id"`
	Index         int32        `json:"index"`
	Name          string       `json:"name"`
	MediaDuration int32        `json:"mediaDuration"`
	Sections      []*SectionVO `json:"sections"`
}

// SectionVO 小节信息及学习进度
type SectionVO struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Index         int32  `json:"index"`
	Type          int32  `json:"type"`
	MediaDuration int32  `json:"mediaDuration"`
	MediaId       int64  `json:"mediaId"`
	Trailer       bool   `json:"trailer"`
	SubjectNum    int32  `json:"subjectNum"`
	HasTest       bool   `json:"hasTest"`
	Moment        int32  `json:"moment"`
	Finished      bool   `json:"finished"`
}

// ---------- 课程分页与列表 ----------

// CoursePageQuery 管理端课程搜索接口查询参数
type CoursePageQuery struct {
	page.Req
	Keyword     string `form:"keyword,optional"`
	FirstCateId int64  `form:"firstCateId,optional"`
	SecondCateId int64 `form:"secondCateId,optional"`
	ThirdCateId int64  `form:"thirdCateId,optional"`
	CourseType  int64  `form:"courseType,optional"`
	Free        string `form:"free,optional"`
	Status      int64  `form:"status"`
	BeginTime   string `form:"beginTime,optional"`
	EndTime     string `form:"endTime,optional"`
}

// CoursePageVO 课程信息（分页）
type CoursePageVO struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	CoverUrl      string `json:"coverUrl"`
	Categories    string `json:"categories"`
	Price         int64  `json:"price"`
	Status        int32  `json:"status"`
	Step          int32  `json:"step"`
	Score         int32  `json:"score"`
	Sold          int32  `json:"sold"`
	Sections      int32  `json:"sections"`
	PublishTime   string `json:"publishTime"`
	PurchaseEndTime string `json:"purchaseEndTime"`
	UpdateTime    string `json:"updateTime"`
	UpdaterName   string `json:"updaterName"`
}

// CourseSimpleInfoDTO 课程简要信息（内部调用）
type CourseSimpleInfoDTO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	CoverUrl        string `json:"coverUrl"`
	Price           int32  `json:"price"`
	Free            bool   `json:"free"`
	Status          int32  `json:"status"`
	ValidDuration   int32  `json:"validDuration"`
	SectionNum      int32  `json:"sectionNum"`
	PurchaseEndTime string `json:"purchaseEndTime"`
	FirstCateId     int64  `json:"firstCateId"`
	SecondCateId    int64  `json:"secondCateId"`
	ThirdCateId     int64  `json:"thirdCateId"`
}

// ---------- 请求查询参数 ----------

// CategoryListQuery 查询课程分类信息
type CategoryListQuery struct {
	Name   string `form:"name,optional"`
	Status int64  `form:"status,optional"`
}

// CourseBaseInfoQuery 获取课程基础信息
type CourseBaseInfoQuery struct {
	See bool `form:"see,optional,default=1"`
}

// CourseCatasQuery 获取课程章节
type CourseCatasQuery struct {
	See          bool `form:"see,optional,default=1"`
	WithPractice bool `form:"withPractice,optional,default=1"`
}

// CourseTeacherQuery 查询课程老师信息
type CourseTeacherQuery struct {
	See bool `form:"see,optional,default=1"`
}

// CheckNameQuery 校验课程名称
type CheckNameQuery struct {
	Id   int64  `form:"id,optional"`
	Name string `form:"name"`
}

// CourseInfoQuery 内部调用-获取课程信息
type CourseInfoQuery struct {
	WithCatalogue bool `form:"withCatalogue,optional"`
	WithTeachers  bool `form:"withTeachers,optional"`
}

// IdsQuery 内部调用-id 列表（重复参数）
type IdsQuery struct {
	Ids []int64 `form:"ids"`
}

// TeacherIdsQuery 内部调用-老师 id 列表
type TeacherIdsQuery struct {
	TeacherIds []int64 `form:"teacherIds"`
}

// MediaIdsQuery 内部调用-媒资 id 列表
type MediaIdsQuery struct {
	MediaIds []int64 `form:"mediaIds"`
}

// SimpleInfoListQuery 课程简要信息列表查询
type SimpleInfoListQuery struct {
	Ids         []int64 `form:"ids,optional"`
	ThirdCataIds []int64 `form:"thirdCataIds,optional"`
}

// NameQuery 内部调用-课程名称
type NameQuery struct {
	Name string `form:"name"`
}

// ---------- 内部调用 ----------

// SubNumAndCourseNumDTO 老师 id 和老师对应的课程数、出题数
type SubNumAndCourseNumDTO struct {
	TeacherId  int64 `json:"teacherId"`
	CourseNum  int32 `json:"courseNum"`
	SubjectNum int32 `json:"subjectNum"`
}

// MediaQuoteDTO 媒资被引用情况
type MediaQuoteDTO struct {
	MediaId  int64 `json:"mediaId"`
	QuoteNum int32 `json:"quoteNum"`
}

// SectionInfoDTO 小节信息，包含课程 id 和媒资 id
type SectionInfoDTO struct {
	CourseId     int64 `json:"courseId"`
	MediaId      int64 `json:"mediaId"`
	Trailer      bool  `json:"trailer"`
	FreeDuration int32 `json:"freeDuration"`
}

// CatalogueDTO 课程章信息（内部调用）
type CatalogueDTO struct {
	Id            int64            `json:"id"`
	Index         int32            `json:"index"`
	Name          string           `json:"name"`
	MediaDuration int32            `json:"mediaDuration"`
	Trailer       bool             `json:"trailer"`
	MediaName     string           `json:"mediaName"`
	MediaId       int64            `json:"mediaId"`
	Type          int32            `json:"type"`
	SubjectNum    int32            `json:"subjectNum"`
	TotalScore    int32            `json:"totalScore"`
	CanUpdate     bool             `json:"canUpdate"`
	Sections      []*CatalogueDTO  `json:"sections"`
}

// CourseFullInfoDTO 课程详细信息，包含课程、目录、教师
type CourseFullInfoDTO struct {
	Id              int64           `json:"id"`
	Name            string          `json:"name"`
	Price           int32           `json:"price"`
	ValidDuration   int32           `json:"validDuration"`
	CoverUrl        string          `json:"coverUrl"`
	PurchaseEndTime string          `json:"purchaseEndTime"`
	FirstCateId     int64           `json:"firstCateId"`
	SecondCateId    int64           `json:"secondCateId"`
	ThirdCateId     int64           `json:"thirdCateId"`
	SectionNum      int32           `json:"sectionNum"`
	Chapters        []*CatalogueDTO `json:"chapters"`
	TeacherIds      []int64         `json:"teacherIds"`
}

// CourseDTO 课程信息（搜索索引库使用）
type CourseDTO struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	CoverUrl      string `json:"coverUrl"`
	Price         int32  `json:"price"`
	Free          bool   `json:"free"`
	Status        string `json:"status"`
	Step          int32  `json:"step"`
	Score         int32  `json:"score"`
	Sold          int32  `json:"sold"`
	Sections      int32  `json:"sections"`
	Duration      int32  `json:"duration"`
	PublishTime   string `json:"publishTime"`
	CreateTime    string `json:"createTime"`
	UpdateTime    string `json:"updateTime"`
	Updater       int64  `json:"updater"`
	Teacher       int64  `json:"teacher"`
	CourseType    int32  `json:"courseType"`
	Enable        int32  `json:"enable"`
	ValidDuration int32  `json:"validDuration"`
	CategoryIdLv1 int64  `json:"categoryIdLv1"`
	CategoryIdLv2 int64  `json:"categoryIdLv2"`
	CategoryIdLv3 int64  `json:"categoryIdLv3"`
}
