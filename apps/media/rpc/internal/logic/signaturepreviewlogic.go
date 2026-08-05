package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignaturePreviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSignaturePreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignaturePreviewLogic {
	return &SignaturePreviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SignaturePreviewLogic) SignaturePreview(in *pb.SignatureRequest) (*pb.SignatureVO, error) {
	// todo: add your logic here and delete this line

	return &pb.SignatureVO{}, nil
}
