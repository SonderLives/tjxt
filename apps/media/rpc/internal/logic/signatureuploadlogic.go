package logic

import (
	"context"
	"time"

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

	// 配置了对象存储：生成预签名 PUT URL（客户端直传，URL 本身即上传凭证）；
	// 未配置：回退 mock 地址并落一条待上传的文件记录
	key := objectKey(in.FileName)
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

	token := ""
	uploadURL := mockUploadURL(key)
	if st := l.svcCtx.Store; st != nil {
		u, err := st.PresignUpload(l.ctx, key, 30*time.Minute)
		if err != nil {
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "生成上传签名失败")
		}
		uploadURL = u
	}

	return &pb.SignatureVO{
		Token:     token,
		Url:       "",
		UploadUrl: uploadURL,
		PlayUrl:   "",
	}, nil
}
