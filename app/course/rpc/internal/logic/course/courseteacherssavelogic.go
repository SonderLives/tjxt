package courselogic

import (
	"context"

	"tjxt/service/course/rpc/internal/model"
	"tjxt/service/course/rpc/internal/svc"
	"tjxt/service/course/rpc/pb"
	"tjxt/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseTeachersSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseTeachersSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeachersSaveLogic {
	return &CourseTeachersSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseTeachersSave 保存课程老师草稿：先清空旧关系，再按入参顺序重新插入。
func (l *CourseTeachersSaveLogic) CourseTeachersSave(in *pb.CourseTeachersSaveRequest) (*pb.Empty, error) {
	if in.Id == 0 {
		return nil, errorx.BadRequestf("课程id不能为空")
	}
	if err := l.svcCtx.CourseTeacherDraftModel.DeleteByCourseId(l.ctx, in.Id); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "清理课程老师关系失败")
	}

	for i, t := range in.Teachers {
		if t == nil {
			continue
		}
		isShow := int64(0)
		if t.IsShow {
			isShow = 1
		}
		data := &model.CourseTeacherDraft{
			CourseId:  in.Id,
			TeacherId: t.Id,
			IsShow:    isShow,
			CIndex:    int64(i + 1),
		}
		if _, err := l.svcCtx.CourseTeacherDraftModel.Insert(l.ctx, data); err != nil {
			return nil, errorx.Wrap(err, errorx.CodeInternal, "保存课程老师关系失败")
		}
	}
	return &pb.Empty{}, nil
}
