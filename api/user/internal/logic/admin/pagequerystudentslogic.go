package admin

import (
	"context"

	"tjxt/apps/user/api/internal/logic/convert"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageQueryStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryStudentsLogic {
	return &PageQueryStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PageQueryStudents 学员分页查询。
func (l *PageQueryStudentsLogic) PageQueryStudents(req *types.UserPageReq) (resp *types.StudentPageVOList, err error) {
	out, err := l.svcCtx.UserRpc.PageQueryStudents(l.ctx, convert.FromUserPageReq(req))
	if err != nil {
		return nil, err
	}
	return convert.ToStudentPageVOList(out), nil
}
