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

type CourseGeneratorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseGeneratorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseGeneratorLogic {
	return &CourseGeneratorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseGenerator 生成课程 id（透传 RPC）。
func (l *CourseGeneratorLogic) CourseGenerator() (resp *types.CourseCataIdVO, err error) {
	res, gerr := l.svcCtx.CourseRpc.CourseGenerator(l.ctx, &pb.Empty{})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "生成课程 id 失败")
	}
	return &types.CourseCataIdVO{Id: res.Id}, nil
}
