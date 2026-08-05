package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"
	"tjxt/pkg/xerr"

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
	if in.MediaId <= 0 {
		return nil, xerr.BadRequestf("mediaId 不能为空")
	}
	media, err := l.svcCtx.MediaModel.FindOneNotDeleted(l.ctx, in.MediaId)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("媒资不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询媒资失败")
	}

	playURL := media.MediaUrl
	if playURL == "" {
		playURL = mockPlayURL(media.Id)
	}
	// mock 实现：真实项目返回 COS/OSS 带签名的播放地址（可追加 token 参数）
	return &pb.SignatureVO{PlayUrl: playURL}, nil
}
