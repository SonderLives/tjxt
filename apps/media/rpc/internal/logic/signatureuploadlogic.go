package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignatureUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSignatureUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignatureUploadLogic {
	return &SignatureUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 签名相关
func (l *SignatureUploadLogic) SignatureUpload(in *pb.SignatureRequest) (*pb.SignatureVO, error) {
	// todo: add your logic here and delete this line

	return &pb.SignatureVO{}, nil
}
