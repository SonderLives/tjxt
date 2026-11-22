// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package signature

import (
	"context"

	backendlogic "tjxt/api/media/backend/logic"
	mediaclient "tjxt/api/media/backend/pb"
	"tjxt/api/media/internal/svc"
	"tjxt/api/media/internal/types"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignaturePreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignaturePreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignaturePreviewLogic {
	return &SignaturePreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignaturePreviewLogic) SignaturePreview(req *types.SignatureReq) (resp *types.SignatureVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := backendlogic.NewSignaturePreviewLogic(l.ctx, l.svcCtx.Backend).SignaturePreview(&mediaclient.SignatureRequest{
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
