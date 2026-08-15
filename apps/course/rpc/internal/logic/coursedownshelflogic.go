package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseDownShelfLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseDownShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDownShelfLogic {
	return &CourseDownShelfLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseDownShelf 单课下架：正式课程状态置为已下架，草稿与子表数据不变。
func (l *CourseDownShelfLogic) CourseDownShelf(in *pb.IdRequest) (*pb.Empty, error) {
	if err := downCourse(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

// downCourse 下架单个课程，CourseDownShelf（单课）与 CourseDown（批量）共用。
func downCourse(ctx context.Context, svcCtx *svc.ServiceContext, id int64) error {
	course, err := svcCtx.CourseModel.FindOne(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return xerr.NotFound("课程不存在")
		}
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}
	if course.Status == int64(CourseStatusDownShelf) {
		return nil
	}
	course.Status = int64(CourseStatusDownShelf)
	if err = svcCtx.CourseModel.Update(ctx, course); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "课程下架失败")
	}
	// 发布下架事件，触发 search 服务从 ES 索引删除（best-effort：失败仅告警）
	svcCtx.PublishCourseEvent(ctx, id, false)
	return nil
}
