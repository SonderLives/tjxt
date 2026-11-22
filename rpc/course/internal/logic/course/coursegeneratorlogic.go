package courselogic

import (
	"context"

	"tjxt/rpc/course/internal/model"
	"tjxt/rpc/course/internal/svc"
	"tjxt/rpc/course/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseGeneratorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseGeneratorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseGeneratorLogic {
	return &CourseGeneratorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseGenerator 预生成一条空课程草稿，返回课程 id 供后续分步保存使用。
func (l *CourseGeneratorLogic) CourseGenerator(in *pb.Empty) (*pb.IdResponse, error) {
	id := nextID()
	draft := &model.CourseDraft{
		Id:        id,
		Status:    CourseStatusPending,
		Step:      0,
		CanUpdate: 1,
		DepId:     0,
		Creater:   0,
		Updater:   0,
	}
	if _, err := l.svcCtx.CourseDraftModel.Insert(l.ctx, draft); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "生成课程草稿失败")
	}
	return &pb.IdResponse{Id: id}, nil
}
