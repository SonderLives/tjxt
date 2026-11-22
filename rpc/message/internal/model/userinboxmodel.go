package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserInboxModel = (*customUserInboxModel)(nil)

type (
	// UserInboxModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInboxModel.
	UserInboxModel interface {
		userInboxModel
		// FindCount 统计某用户未过期的站内信数量（可选类型过滤）
		FindCount(ctx context.Context, userId, msgType int64) (int64, error)
		// FindList 分页查询某用户未过期的站内信，按 push_time 倒序
		FindList(ctx context.Context, userId, msgType int64, offset, limit int64) ([]*UserInbox, error)
	}

	customUserInboxModel struct {
		*defaultUserInboxModel
	}
)

// NewUserInboxModel returns a model for the database table.
func NewUserInboxModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserInboxModel {
	return &customUserInboxModel{
		defaultUserInboxModel: newUserInboxModel(conn, c, opts...),
	}
}

// FindCount 统计某用户未过期的站内信数量，msgType <= 0 时不做类型过滤
func (m *customUserInboxModel) FindCount(ctx context.Context, userId, msgType int64) (int64, error) {
	where, args := inboxWhere(userId, msgType)
	var total int64
	query := fmt.Sprintf("select count(1) from %s%s", m.table, where)
	if err := m.QueryRowNoCacheCtx(ctx, &total, query, args...); err != nil {
		return 0, err
	}
	return total, nil
}

// FindList 分页查询某用户未过期的站内信，msgType <= 0 时不做类型过滤
func (m *customUserInboxModel) FindList(ctx context.Context, userId, msgType int64, offset, limit int64) ([]*UserInbox, error) {
	where, args := inboxWhere(userId, msgType)
	args = append(args, limit, offset)
	var list []*UserInbox
	query := fmt.Sprintf("select %s from %s%s order by `push_time` desc limit ? offset ?", userInboxRows, m.table, where)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func inboxWhere(userId, msgType int64) (string, []any) {
	cond := []string{"`user_id` = ?", "`expire_time` > NOW()"}
	args := []any{userId}
	if msgType > 0 {
		cond = append(cond, "`type` = ?")
		args = append(args, msgType)
	}
	return " where " + strings.Join(cond, " and "), args
}
