package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NoticeTaskTargetModel = (*customNoticeTaskTargetModel)(nil)

type (
	// NoticeTaskTargetModel 通知任务目标用户关联表的轻量查询模型（无缓存、无增删改）。
	NoticeTaskTargetModel interface {
		// FindTargetIdsByTaskId 查询任务的目标用户 id 列表
		FindTargetIdsByTaskId(ctx context.Context, taskId int64) ([]int64, error)
	}

	customNoticeTaskTargetModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewNoticeTaskTargetModel 创建轻量模型，直接使用 sqlx 连接查询。
func NewNoticeTaskTargetModel(conn sqlx.SqlConn) NoticeTaskTargetModel {
	return &customNoticeTaskTargetModel{
		conn:  conn,
		table: "`notice_task_target`",
	}
}

// FindTargetIdsByTaskId 查询任务的目标用户 id 列表
func (m *customNoticeTaskTargetModel) FindTargetIdsByTaskId(ctx context.Context, taskId int64) ([]int64, error) {
	var list []int64
	query := fmt.Sprintf("select `target_id` from %s where `task_id` = ?", m.table)
	if err := m.conn.QueryRowsCtx(ctx, &list, query, taskId); err != nil {
		return nil, err
	}
	return list, nil
}
