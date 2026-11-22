package medialogic

import (
	"context"
	"errors"

	"tjxt/pkg/errorx"
	"tjxt/rpc/media/internal/model"
	"tjxt/rpc/media/internal/svc"
	"tjxt/rpc/media/pb"

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
			if errors.Is(err, model.ErrNotFound) {
				return nil, errorx.NotFound("媒资不存在")
			}
			return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询媒资失败")
		}
		if st := l.svcCtx.Store; st != nil {
			if file, ferr := l.svcCtx.FileModel.FindByKey(l.ctx, media.FileId); ferr == nil {
				if u, perr := st.PresignGet(l.ctx, file.Key, PresignViewTTL); perr == nil {
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
		return nil, errorx.BadRequestf("fileName 不能为空")
	}
	key := objectKey(in.FileName)
	return &pb.SignatureVO{PlayUrl: mockPlayURL(key)}, nil
}
