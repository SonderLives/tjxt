package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	// userDetailColumns 用户详情表列列表
	userDetailColumns = `id, type, name, gender, icon, email, qq, birthday, job, province, city,
	 district, intro, photo, role_id, course_amount, create_time, update_time, creater, updater, dep_id`

	// userDetailScanColumns 扫描用列：可空列以 COALESCE 兜底
	userDetailScanColumns = `id, type, COALESCE(name, ''), gender, COALESCE(icon, ''), email, qq, birthday,
	 COALESCE(job, ''), province, city, district, COALESCE(intro, ''), COALESCE(photo, ''), role_id,
	 COALESCE(course_amount, 0), create_time, update_time, creater, updater, dep_id`

	// Cache key 模式：cache:detail:id:<id>
	cacheDetailKeyPrefix = "cache:detail:id:"
)

// UserDetail 用户详情（user_detail 表）
type UserDetail struct {
	Id           int64
	Type         int64
	Name         string
	Gender       int64
	Icon         string
	Email        string
	Qq           string
	Birthday     sql.NullTime
	Job          string
	Province     string
	City         string
	District     string
	Intro        string
	Photo        string
	RoleId       int64
	CourseAmount int64
	CreateTime   time.Time
	UpdateTime   time.Time
	Creater      sql.NullInt64
	Updater      sql.NullInt64
	DepId        int64
}

// UserDetailModel 用户详情数据访问
type UserDetailModel struct {
	conn  sqlx.SqlConn
	c     cache.Cache
	table string
}

// NewUserDetailModel 创建用户详情数据访问对象，自动接入 Redis 缓存
func NewUserDetailModel(conn sqlx.SqlConn, c cache.Cache) *UserDetailModel {
	return &UserDetailModel{
		conn:  conn,
		c:     c,
		table: "user_detail",
	}
}

// FindById 按关联用户 id 查询详情（带缓存）
// 缓存 key 模式：cache:detail:id:<id>
func (m *UserDetailModel) FindById(ctx context.Context, id int64) (*UserDetail, error) {
	cacheKey := fmt.Sprintf("%s%d", cacheDetailKeyPrefix, id)
	var d UserDetail
	if err := m.c.GetCtx(ctx, cacheKey, &d); err == nil {
		return &d, nil
	}

	err := m.conn.QueryRowCtx(ctx, &d, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ?`, userDetailScanColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err == nil && m.c != nil {
		_ = m.c.SetWithExpireCtx(ctx, cacheKey, d, 10*time.Minute)
	}
	return &d, err
}

// Insert 新增详情（清除缓存）
func (m *UserDetailModel) Insert(ctx context.Context, d *UserDetail) error {
	if m.c != nil {
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%d", cacheDetailKeyPrefix, d.Id))
	}
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, type, name, gender, icon, email, qq, birthday, job, province, city,
			district, intro, photo, role_id, course_amount, create_time, update_time)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		m.table),
		d.Id, d.Type, d.Name, d.Gender, d.Icon, d.Email, d.Qq, d.Birthday, d.Job, d.Province, d.City,
		d.District, d.Intro, d.Photo, d.RoleId, d.CourseAmount)
	return err
}

// InsertTx 事务内新增详情（不触发缓存清除）
func (m *UserDetailModel) InsertTx(ctx context.Context, tx sqlx.Session, d *UserDetail) error {
	_, err := tx.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, type, name, gender, icon, email, qq, birthday, job, province, city,
			district, intro, photo, role_id, course_amount, create_time, update_time)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		m.table),
		d.Id, d.Type, d.Name, d.Gender, d.Icon, d.Email, d.Qq, d.Birthday, d.Job, d.Province, d.City,
		d.District, d.Intro, d.Photo, d.RoleId, d.CourseAmount)
	return err
}

// UpdateDetail 更新用户详情（清除缓存）
func (m *UserDetailModel) UpdateDetail(ctx context.Context, d *UserDetail) error {
	if m.c != nil {
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%d", cacheDetailKeyPrefix, d.Id))
	}
	sets := []string{"name = ?", "role_id = ?", "update_time = NOW()"}
	args := []any{d.Name, d.RoleId}
	updates := map[string]any{
		"gender":   d.Gender,
		"icon":     d.Icon,
		"email":    d.Email,
		"qq":       d.Qq,
		"job":      d.Job,
		"province": d.Province,
		"city":     d.City,
		"district": d.District,
		"intro":    d.Intro,
		"photo":    d.Photo,
	}
	for k, v := range updates {
		sets = append(sets, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	args = append(args, d.Id)
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET %s WHERE id = ?`, m.table, strings.Join(sets, ", ")), args...)
	return err
}

// UpdateCourseAmount 更新学员已购课程数（清除缓存）
func (m *UserDetailModel) UpdateCourseAmount(ctx context.Context, id, amount int64) error {
	if m.c != nil {
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%d", cacheDetailKeyPrefix, id))
	}
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET course_amount = ?, update_time = NOW() WHERE id = ?`, m.table), amount, id)
	return err
}

// PageCond 分页查询条件
type PageCond struct {
	UserType int64
	Name     string
	Phone    string
	Status   int64
	Offset   int64
	Limit    int64
	IsAsc    bool
	SortBy   string // 排序字段：id | create_time（由调用方校验合法性）
}

// PageRow 分页查询结果行
type PageRow struct {
	Id           int64
	Username     string
	CellPhone    string
	Type         int64
	Status       int64
	Name         string
	Gender       int64
	Icon         string
	Photo        string
	Job          string
	Intro        string
	RoleId       int64
	CourseAmount int64
	CreateTime   time.Time
}

// ListPage 分页查询指定类型的用户信息（联表 user + user_detail）
func (m *UserDetailModel) ListPage(ctx context.Context, cond *PageCond) ([]PageRow, int64, error) {
	where := []string{"u.type = ?"}
	args := []any{cond.UserType}
	if cond.Name != "" {
		where = append(where, "d.name LIKE ?")
		args = append(args, "%"+cond.Name+"%")
	}
	if cond.Phone != "" {
		where = append(where, "u.cell_phone LIKE ?")
		args = append(args, "%"+cond.Phone+"%")
	}
	if cond.Status > 0 {
		where = append(where, "u.status = ?")
		args = append(args, cond.Status)
	}
	whereSQL := strings.Join(where, " AND ")

	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(*) FROM %s d JOIN `+"`user`"+` u ON d.id = u.id WHERE %s`, m.table, whereSQL), args...)
	if err != nil {
		return nil, 0, err
	}

	// 默认按 create_time 降序
	sortBy := "create_time"
	if cond.SortBy != "" {
		sortBy = cond.SortBy
	}
	order := sortBy + " DESC"
	if cond.IsAsc {
		order = sortBy + " ASC"
	}
	args = append(args, cond.Limit, cond.Offset)
	rows := make([]PageRow, 0, cond.Limit)
	err = m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT u.id, u.username, u.cell_phone, u.type, u.status,
		        COALESCE(d.name, ''), d.gender, COALESCE(d.icon, ''), COALESCE(d.photo, ''),
		        COALESCE(d.job, ''), COALESCE(d.intro, ''), d.role_id, COALESCE(d.course_amount, 0),
		        u.create_time
		 FROM %s d JOIN `+"`user`"+` u ON d.id = u.id
		 WHERE %s ORDER BY %s LIMIT ? OFFSET ?`, m.table, whereSQL, order), args...)
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListByIDs 按 id 列表批量查询详情（无缓存）
func (m *UserDetailModel) ListByIDs(ctx context.Context, ids []int64) (map[int64]*UserDetail, error) {
	result := make(map[int64]*UserDetail, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var rows []UserDetail
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id IN (%s)`, userDetailScanColumns, m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		result[rows[i].Id] = &rows[i]
	}
	return result, nil
}