package student

import (
	"context"

	"tjxt/apps/user/api/internal/logic/convert"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStudentPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateStudentPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStudentPasswordLogic {
	return &UpdateStudentPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateStudentPassword 学员修改密码（公开接口），转发至 user RPC。
func (l *UpdateStudentPasswordLogic) UpdateStudentPassword(req *types.StudentFormReq) error {
	_, err := l.svcCtx.UserRpc.UpdateStudentPassword(l.ctx, convert.FromStudentFormReq(req))
	return err
}
