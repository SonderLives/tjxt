package logic

import (
	"context"
	"time"

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
	// 指定媒资：配置对象存储时按 file.key 出预签名预览地址，否则回退 mock
	if in.MediaId > 0 {
		media, err := l.svcCtx.MediaModel.FindOneNotDeleted(l.ctx, in.MediaId)
		if err != nil {
			if isNotFound(err) {
				return nil, xerr.NotFound("媒资不存在")
			}
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询媒资失败")
		}
		if st := l.svcCtx.Store; st != nil {
			if file, ferr := l.svcCtx.FileModel.FindByKey(l.ctx, media.FileId); ferr == nil {
				if u, perr := st.PresignGet(l.ctx, file.Key, 2*time.Hour); perr == nil {
					return &pb.SignatureVO{PlayUrl: u}, nil
				} else {
					logx.WithContext(l.ctx).Errorf("presign preview url failed, mediaId=%d: %v", in.MediaId, perr)
				}
			}
		}
		return &pb.SignatureVO{PlayUrl: mockPlayURL(in.MediaId)}, nil
	}

	// 未指定媒资：按文件名生成占位 key 返回
	if in.FileName == "" {
		return nil, xerr.BadRequestf("fileName 不能为空")
	}
	key := objectKey(in.FileName)
	return &pb.SignatureVO{PlayUrl: mockPlayURL(key)}, nil
}
