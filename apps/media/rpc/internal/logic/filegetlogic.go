package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileGetLogic {
	return &FileGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 文件管理
func (l *FileGetLogic) FileGet(in *pb.FileIdRequest) (*pb.FileVO, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("id 不能为空")
	}
	file, err := l.svcCtx.FileModel.FindOneNotDeleted(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("文件不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询文件失败")
	}
	return toFileVO(file), nil
}
