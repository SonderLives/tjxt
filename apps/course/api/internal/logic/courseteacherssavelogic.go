// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseTeachersSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTeachersSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeachersSaveLogic {
	return &CourseTeachersSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseTeachersSave 保存课程老师（透传 RPC）。
func (l *CourseTeachersSaveLogic) CourseTeachersSave(req *types.CourseTeacherSaveReq) (resp *types.NameExistVO, err error) {
	teachers := make([]*pb.TeacherShowBind, 0, len(req.Teachers))
	for _, t := range req.Teachers {
		teachers = append(teachers, &pb.TeacherShowBind{
			Id:     t.Id,
			IsShow: t.IsShow,
		})
	}
	_, err = l.svcCtx.CourseRpc.CourseTeachersSave(l.ctx, &pb.CourseTeachersSaveRequest{
		Id:       req.Id,
		Teachers: teachers,
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "保存课程老师失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
