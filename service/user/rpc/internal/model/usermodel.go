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

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		// FindPageByType 按用户类型 + 条件分页查询，LEFT JOIN user_detail 返回联合视图。
		// status < 0 表示不过滤状态；sortCol/sortDir 由调用方白名单映射后传入，禁止透传原始用户输入。
		FindPageByType(ctx context.Context, userType int64, name, phone string, status, offset, limit int64, sortCol, sortDir string) ([]*UserWithDetail, int64, error)
		// FindByIdsWithDetail 批量查询用户并附带资料，用于 GetUsersByIds。
		FindByIdsWithDetail(ctx context.Context, ids []int64) ([]*UserWithDetail, error)
		// ExistsByCellPhone 判断手机号是否已被任意类型用户占用。
		ExistsByCellPhone(ctx context.Context, cellPhone string) (bool, error)
		// UpdateStatus 更新账户状态并失效相关缓存键。
		UpdateStatus(ctx context.Context, id, status, updater int64) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// userWithDetailRows 联合视图的投影列，列名需与 UserWithDetail 的 db tag 对应。
const userWithDetailRows = "u.`id`, u.`cell_phone`, u.`type`, u.`status`, u.`create_time`, " +
	"d.`name`, d.`icon`, d.`gender`, d.`photo`, d.`job`, d.`intro`, d.`course_amount`, d.`role_id`, " +
	"d.`email`, d.`qq`, d.`province`, d.`city`, d.`district`"

// UserWithDetail 用户与其资料（user_detail）的联合视图，用于详情与分页查询。
type UserWithDetail struct {
	Id           int64          `db:"id"`
	CellPhone    string         `db:"cell_phone"`
	Type         int64          `db:"type"`
	Status       int64          `db:"status"`
	CreateTime   time.Time      `db:"create_time"`
	Name         sql.NullString `db:"name"`
	Icon         sql.NullString `db:"icon"`
	Gender       int64          `db:"gender"`
	Photo        sql.NullString `db:"photo"`
	Job          sql.NullString `db:"job"`
	Intro        sql.NullString `db:"intro"`
	CourseAmount int64          `db:"course_amount"`
	RoleId       int64          `db:"role_id"`
	Email        sql.NullString `db:"email"`
	Qq           sql.NullString `db:"qq"`
	Province     sql.NullString `db:"province"`
	City         sql.NullString `db:"city"`
	District     sql.NullString `db:"district"`
}

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

// FindPageByType 按类型与条件分页查询，联合 user_detail。
func (m *customUserModel) FindPageByType(ctx context.Context, userType int64, name, phone string, status, offset, limit int64, sortCol, sortDir string) ([]*UserWithDetail, int64, error) {
	where := []string{"u.`type` = ?"}
	args := []any{userType}
	if name != "" {
		where = append(where, "d.`name` like ?")
		args = append(args, "%"+name+"%")
	}
	if phone != "" {
		where = append(where, "u.`cell_phone` like ?")
		args = append(args, "%"+phone+"%")
	}
	if status >= 0 {
		where = append(where, "u.`status` = ?")
		args = append(args, status)
	}
	whereClause := strings.Join(where, " AND ")
	base := "FROM `user` u LEFT JOIN `user_detail` d ON u.`id` = d.`id` WHERE " + whereClause
	query := fmt.Sprintf("SELECT %s %s ORDER BY %s %s LIMIT ?, ?", userWithDetailRows, base, sortCol, sortDir)
	countQuery := "SELECT count(*) " + base

	var (
		list  []*UserWithDetail
		total int64
	)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, append(append([]any{}, args...), offset, limit)...); err != nil {
		return nil, 0, err
	}
	if err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// FindByIdsWithDetail 批量查询用户并附带资料。
func (m *customUserModel) FindByIdsWithDetail(ctx context.Context, ids []int64) ([]*UserWithDetail, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	holders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf("SELECT %s FROM `user` u LEFT JOIN `user_detail` d ON u.`id` = d.`id` WHERE u.`id` IN (%s) ORDER BY u.`id` DESC", userWithDetailRows, holders)
	var list []*UserWithDetail
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...)
	return list, err
}

// ExistsByCellPhone 判断手机号是否已被占用。
func (m *customUserModel) ExistsByCellPhone(ctx context.Context, cellPhone string) (bool, error) {
	query := "SELECT count(*) FROM `user` WHERE `cell_phone` = ?"
	var cnt int64
	if err := m.QueryRowNoCacheCtx(ctx, &cnt, query, cellPhone); err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// UpdateStatus 更新账户状态并失效缓存。
func (m *customUserModel) UpdateStatus(ctx context.Context, id, status, updater int64) error {
	u, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}
	query := "update `user` set `status` = ?, `updater` = ?, `update_time` = now() where `id` = ?"
	if _, err := m.ExecNoCacheCtx(ctx, query, status, updater, id); err != nil {
		return err
	}
	_ = m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheTjUserUserIdPrefix, id))
	_ = m.DelCacheCtx(ctx, fmt.Sprintf("%s%v:%v", cacheTjUserUserCellPhoneTypePrefix, u.CellPhone, u.Type))
	_ = m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheTjUserUserUsernamePrefix, u.Username))
	return nil
}
