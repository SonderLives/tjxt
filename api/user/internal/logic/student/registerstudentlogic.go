package student

import (
	"context"

	"tjxt/api/user/internal/logic/convert"
	"tjxt/api/user/internal/svc"
	"tjxt/api/user/internal/types"

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

// RegisterStudent 学员注册（公开接口），转发至 user RPC。
func (l *RegisterStudentLogic) RegisterStudent(req *types.StudentFormReq) error {
	_, err := l.svcCtx.UserRpc.RegisterStudent(l.ctx, convert.FromStudentFormReq(req))
	return err
}
