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

type CourseSubjectsSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSubjectsSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSubjectsSaveLogic {
	return &CourseSubjectsSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseSubjectsSave 保存课程目录题目（透传 RPC）。
func (l *CourseSubjectsSaveLogic) CourseSubjectsSave(req *types.CourseSubjectsSaveReq) (resp *types.NameExistVO, err error) {
	subjects := make([]*pb.CatalogSubjectBind, 0, len(req.List))
	for _, item := range req.List {
		subjects = append(subjects, &pb.CatalogSubjectBind{
			CataId:     item.CataId,
			SubjectIds: item.SubjectIds,
		})
	}
	_, err = l.svcCtx.CourseRpc.CourseSubjectsSave(l.ctx, &pb.CourseSubjectsSaveRequest{
		CourseId: req.Id,
		Subjects: subjects,
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "保存课程题目失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
