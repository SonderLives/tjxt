package userlogic

import (
	"context"

	"tjxt/app/user/rpc/internal/svc"
	"tjxt/app/user/rpc/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckCellPhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckCellPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckCellPhoneLogic {
	return &CheckCellPhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CheckCellPhone 检查手机号是否已被注册（true=已存在）。
func (l *CheckCellPhoneLogic) CheckCellPhone(in *pb.CheckCellPhoneRequest) (*pb.BoolResponse, error) {
	exists, err := l.svcCtx.UserModel.ExistsByCellPhone(l.ctx, in.CellPhone)
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "校验手机号失败")
	}
	return &pb.BoolResponse{Result: exists}, nil
}
