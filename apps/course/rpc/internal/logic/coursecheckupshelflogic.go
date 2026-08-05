package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// courseStepFinished 课程信息填写完成的进度值：
// 1 基本信息、2 课程目录、3 课程视频、4 课程题目、5 课程老师，达到 5 才允许上架。
const courseStepFinished int64 = 5

type CourseCheckUpShelfLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCheckUpShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCheckUpShelfLogic {
	return &CourseCheckUpShelfLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseCheckUpShelf 上架前校验：课程草稿的填写进度必须达到 5（信息/目录/视频/题目/老师均已保存）。
// proto 的返回值是 Empty（无 existed 字段），因此校验结果通过错误传递：
// 返回 nil 表示可上架（API 侧对应 existed=true），返回 BadRequest 错误表示不可上架（existed=false）。
func (l *CourseCheckUpShelfLogic) CourseCheckUpShelf(in *pb.IdRequest) (*pb.Empty, error) {
	draft, err := l.svcCtx.CourseDraftModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("课程草稿不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	if draft.Step < courseStepFinished {
		return nil, xerr.BadRequestf("课程信息未填写完整，当前进度 %d/%d，不能上架", draft.Step, courseStepFinished)
	}
	return &pb.Empty{}, nil
}
