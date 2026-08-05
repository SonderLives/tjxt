package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/rpc/internal/svc"
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
func (l *PlanPageLogic) PlanPage(in *pb.LessonPageRequest) (*pb.PlanPageReply, error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := l.svcCtx.LearningService.ListLessonPlans(l.ctx, userID, in.PageNo, in.PageSize, in.IsAsc)
	if err != nil {
		return nil, err
	}
	var weekTotalPlan int64
	vos := make([]*pb.LearningLessonVO, 0, len(list))
	for _, lsn := range list {
		vos = append(vos, toLessonVO(lsn))
		weekTotalPlan += nullInt64(lsn.WeekFreq)
	}
	return &pb.PlanPageReply{
		Total:         total,
		Pages:         calcPages(total, in.PageSize),
		WeekTotalPlan: int32(weekTotalPlan),
		// WeekFinished / WeekPoints 暂无独立数据源，留 0
		List: vos,
	}, nil
}
