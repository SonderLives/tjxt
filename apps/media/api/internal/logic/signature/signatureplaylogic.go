// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package signature

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignaturePlayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignaturePlayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignaturePlayLogic {
	return &SignaturePlayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignaturePlayLogic) SignaturePlay(req *types.SignatureReq) (resp *types.SignatureVO, err error) {
	// todo: add your logic here and delete this line

	return
}
