// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package signature

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"
	mediaclient "tjxt/apps/media/rpc/media"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignatureUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignatureUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignatureUploadLogic {
	return &SignatureUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignatureUploadLogic) SignatureUpload(req *types.SignatureReq) (resp *types.SignatureVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.MediaRpc.SignatureUpload(l.ctx, &mediaclient.SignatureRequest{
		MediaId:   req.MediaId,
		FileName:  req.FileName,
		MediaType: req.MediaType,
	})
	if err != nil {
		return nil, err
	}
	return &types.SignatureVO{
		Token:     rpcResp.Token,
		Url:       rpcResp.Url,
		UploadUrl: rpcResp.UploadUrl,
		PlayUrl:   rpcResp.PlayUrl,
	}, nil
}
