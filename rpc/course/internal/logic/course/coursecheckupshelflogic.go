package courselogic

import (
	"context"

	"tjxt/rpc/course/internal/svc"
	"tjxt/rpc/course/pb"
	"tjxt/pkg/errorx"

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
			return nil, errorx.NotFound("课程草稿不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询课程草稿失败")
	}
	if draft.Step < courseStepFinished {
		return nil, errorx.BadRequestf("课程信息未填写完整，当前进度 %d/%d，不能上架", draft.Step, courseStepFinished)
	}
	return &pb.Empty{}, nil
}
