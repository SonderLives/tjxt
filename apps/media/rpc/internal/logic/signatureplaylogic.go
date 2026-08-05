package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignaturePlayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSignaturePlayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignaturePlayLogic {
	return &SignaturePlayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SignaturePlayLogic) SignaturePlay(in *pb.SignatureRequest) (*pb.SignatureVO, error) {
	// todo: add your logic here and delete this line

	return &pb.SignatureVO{}, nil
}
