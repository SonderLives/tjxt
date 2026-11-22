package admin

import (
	"context"

	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	pb "tjxt/apps/user/rpc/pb"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ResetPassword 重置用户密码为默认口令。
func (l *ResetPasswordLogic) ResetPassword(req *types.IdPathReq) error {
	_, err := l.svcCtx.UserRpc.ResetPassword(l.ctx, &pb.UserIdRequest{UserId: req.Id})
	return err
}
