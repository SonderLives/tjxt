package admin

import (
	"context"

	"tjxt/api/user/internal/logic/convert"
	"tjxt/api/user/internal/svc"
	"tjxt/api/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryTeachersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageQueryTeachersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryTeachersLogic {
	return &PageQueryTeachersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PageQueryTeachers 老师分页查询。
func (l *PageQueryTeachersLogic) PageQueryTeachers(req *types.UserPageReq) (resp *types.TeacherPageVOList, err error) {
	out, err := l.svcCtx.UserRpc.PageQueryTeachers(l.ctx, convert.FromUserPageReq(req))
	if err != nil {
		return nil, err
	}
	return convert.ToTeacherPageVOList(out), nil
}
