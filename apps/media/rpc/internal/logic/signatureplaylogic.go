package logic

import (
	"context"
	"time"

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
	// 配置了对象存储：按 file.key 生成带时效的预签名播放地址（防盗链）；
	// 失败回退媒资持久 URL
	if st := l.svcCtx.Store; st != nil {
		if file, ferr := l.svcCtx.FileModel.FindByKey(l.ctx, media.FileId); ferr == nil {
			if u, perr := st.PresignGet(l.ctx, file.Key, 2*time.Hour); perr == nil {
				playURL = u
			} else {
				logx.WithContext(l.ctx).Errorf("presign play url failed, mediaId=%d: %v", media.Id, perr)
			}
		} else {
			logx.WithContext(l.ctx).Errorf("query file for play url failed, mediaId=%d fileId=%s: %v", media.Id, media.FileId, ferr)
		}
	}
	return &pb.SignatureVO{PlayUrl: playURL}, nil
}
