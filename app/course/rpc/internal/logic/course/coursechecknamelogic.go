package courselogic

import (
	"context"
	"errors"

	"tjxt/common/errorx"
	"tjxt/service/course/rpc/internal/model"
	"tjxt/service/course/rpc/internal/svc"
	"tjxt/service/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCheckNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCheckNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCheckNameLogic {
	return &CourseCheckNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseCheckName 校验课程名称是否已存在（编辑模式下排除自身）。
func (l *CourseCheckNameLogic) CourseCheckName(in *pb.CourseCheckNameRequest) (*pb.NameExistReply, error) {
	if in.Name == "" {
		return nil, errorx.BadRequestf("课程名称不能为空")
	}
	_, err := l.svcCtx.CourseDraftModel.FindByNameExceptId(l.ctx, in.Name, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &pb.NameExistReply{Existed: false}, nil
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "校验课程名称失败")
	}
	return &pb.NameExistReply{Existed: true}, nil
}
