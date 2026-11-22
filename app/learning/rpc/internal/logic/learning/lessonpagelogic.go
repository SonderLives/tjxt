package learninglogic

import (
	"context"

	"tjxt/pkg/auth"
	"tjxt/pkg/page"

	"tjxt/app/learning/rpc/internal/svc"
	"tjxt/app/learning/rpc/pb"

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
func (l *LessonPageLogic) LessonPage(in *pb.LessonPageRequest) (*pb.LessonPageReply, error) {
	userID, err := auth.UserIdFromIncomingMD(l.ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := l.svcCtx.LearningService.ListLessons(l.ctx, userID, in.PageNo, in.PageSize, in.IsAsc)
	if err != nil {
		return nil, err
	}
	vos := make([]*pb.LearningLessonVO, 0, len(list))
	for _, lsn := range list {
		vos = append(vos, toLessonVO(lsn))
	}
	return &pb.LessonPageReply{
		Total: total,
		Pages: page.CalcPages(total, in.PageSize),
		List:  vos,
	}, nil
}
