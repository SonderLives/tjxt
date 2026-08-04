// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"

	"learning/internal/svc"
	"learning/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePlanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePlanLogic {
	return &CreatePlanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePlanLogic) CreatePlan(req *types.PlanRequest) (resp *types.Result, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	if req.CourseId <= 0 || req.Freq <= 0 {
		return nil, fmt.Errorf("courseId and freq must be positive")
	}
	if err = l.svcCtx.LessonService.CreatePlan(l.ctx, userID, req.CourseId, req.Freq); err != nil {
		return nil, err
	}
	return success(nil), nil
}
