package medialogic

import (
	"context"
	"errors"

	"tjxt/service/media/api/backend/model"
	"tjxt/service/media/api/backend/pb"
	"tjxt/service/media/api/backend/svc"
	"tjxt/common/errorx"

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
		return nil, errorx.BadRequestf("id 不能为空")
	}
	file, err := l.svcCtx.FileModel.FindOneNotDeleted(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("文件不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询文件失败")
	}
	return toFileVO(file, l.svcCtx.Store), nil
}
