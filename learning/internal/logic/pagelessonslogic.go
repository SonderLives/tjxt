// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"learning/internal/svc"
	"learning/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageLessonsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageLessonsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLessonsLogic {
	return &PageLessonsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLessonsLogic) PageLessons(req *types.PageRequest) (resp *types.Result, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	pageNo, pageSize := normalizePage(req.PageNo, req.PageSize)
	lessons, total, err := l.svcCtx.LessonService.ListLessons(l.ctx, userID, pageNo, pageSize)
	if err != nil {
		return nil, err
	}
	items := make([]types.Lesson, 0, len(lessons))
	for index := range lessons {
		items = append(items, lessonResponse(&lessons[index]))
	}
	pages := (total + pageSize - 1) / pageSize
	return success(types.LessonPage{List: items, Total: total, Pages: pages}), nil
}
