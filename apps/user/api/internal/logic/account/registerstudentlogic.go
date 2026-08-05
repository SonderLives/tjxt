// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"context"

	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterStudentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterStudentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterStudentLogic {
	return &RegisterStudentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterStudentLogic) RegisterStudent(req *types.RegisterReq) (resp *types.RegisterVO, err error) {
	// todo: add your logic here and delete this line

	return
}
