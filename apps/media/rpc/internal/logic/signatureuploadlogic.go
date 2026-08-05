package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/model"
	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"
	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/xerr"

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
	if in.FileName == "" {
		return nil, xerr.BadRequestf("fileName 不能为空")
	}
	if !supportedMediaType(in.MediaType) {
		return nil, xerr.BadRequestf("mediaType 仅支持 video/image/audio")
	}

	// mock 实现：真实项目调用腾讯云 COS / 阿里 OSS 的 STS 接口生成临时签名
	// （含临时密钥、签名串与上传 URL），此处仅占位并落一条待上传的文件记录
	key := mockObjectKey(in.FileName)
	file := &model.File{
		Id:       idgen.NextID(),
		Key:      key,
		Filename: in.FileName,
		Status:   FileStatusPending,
		Platform: PlatformTencent,
	}
	if _, err := l.svcCtx.FileModel.Insert(l.ctx, file); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "创建文件记录失败")
	}

	return &pb.SignatureVO{
		Token:     mockToken(),
		Url:       "",
		UploadUrl: mockUploadURL(key),
		PlayUrl:   "",
	}, nil
}
