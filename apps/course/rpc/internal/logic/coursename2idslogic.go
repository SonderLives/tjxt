package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseName2IdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseName2IdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseName2IdsLogic {
	return &CourseName2IdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseName2Ids 按课程名称模糊匹配，返回命中的课程 id 列表。
func (l *CourseName2IdsLogic) CourseName2Ids(in *pb.CourseNameRequest) (*pb.CourseIdList, error) {
	if in.Name == "" {
		return &pb.CourseIdList{Ids: []int64{}}, nil
	}
	list, err := l.svcCtx.CourseModel.FindByNameLike(l.ctx, in.Name)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "按名称查询课程失败")
	}
	ids := make([]int64, 0, len(list))
	for _, c := range list {
		ids = append(ids, c.Id)
	}
	return &pb.CourseIdList{Ids: ids}, nil
}
