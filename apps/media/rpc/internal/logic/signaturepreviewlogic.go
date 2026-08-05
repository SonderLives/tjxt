package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"
	"tjxt/pkg/xerr"

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
	// 指定媒资：校验存在后返回 mock 预览地址
	if in.MediaId > 0 {
		if _, err := l.svcCtx.MediaModel.FindOneNotDeleted(l.ctx, in.MediaId); err != nil {
			if isNotFound(err) {
				return nil, xerr.NotFound("媒资不存在")
			}
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询媒资失败")
		}
		// mock 实现：真实项目返回 COS/OSS 带签名的预览地址
		return &pb.SignatureVO{PlayUrl: mockPlayURL(in.MediaId)}, nil
	}

	// 未指定媒资：按文件名生成占位 key 返回
	if in.FileName == "" {
		return nil, xerr.BadRequestf("fileName 不能为空")
	}
	key := mockObjectKey(in.FileName)
	return &pb.SignatureVO{PlayUrl: mockPlayURL(key)}, nil
}
