package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLessonPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonPageLogic {
	return &LessonPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 我的课表分页
func (l *LessonPageLogic) LessonPage(req *types.PageRequest) (*types.LessonPageReply, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	reply, err := l.svcCtx.LearningRpc.LessonPage(l.ctx, &pb.LessonPageRequest{
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
	return &types.LessonPageReply{
		Total: reply.Total,
		Pages: reply.Pages,
		List:  list,
	}, nil
}
