package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PlanPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPlanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlanPageLogic {
	return &PlanPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 我的学习计划分页
func (l *PlanPageLogic) PlanPage(req *types.PageRequest) (*types.PlanPageReply, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	reply, err := l.svcCtx.LearningRpc.PlanPage(l.ctx, &pb.LessonPageRequest{
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		IsAsc:    req.IsAsc,
		SortBy:   req.SortBy,
	})
	if err != nil {
		return nil, err
	}
	enrichLessons(l.ctx, l.svcCtx, reply.List)
	list := make([]types.LearningLessonVO, 0, len(reply.List))
	for _, v := range reply.List {
		list = append(list, toLessonVOTypes(v))
	}
	return &types.PlanPageReply{
		Total:         reply.Total,
		Pages:         reply.Pages,
		WeekFinished:  int64(reply.WeekFinished),
		WeekPoints:    int64(reply.WeekPoints),
		WeekTotalPlan: int64(reply.WeekTotalPlan),
		List:          list,
	}, nil
}
