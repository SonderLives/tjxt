package student

import (
	"context"

	"tjxt/apps/user/api/internal/logic/convert"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckCellPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckCellPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckCellPhoneLogic {
	return &CheckCellPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CheckCellPhone 检查手机号是否已注册（true=已存在）。
func (l *CheckCellPhoneLogic) CheckCellPhone(req *types.CheckCellPhoneReq) (resp *types.BoolVO, err error) {
	out, err := l.svcCtx.UserRpc.CheckCellPhone(l.ctx, convert.FromCheckCellPhoneReq(req))
	if err != nil {
		return nil, err
	}
	return convert.ToBoolVO(out), nil
}
