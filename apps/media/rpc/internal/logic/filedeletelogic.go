package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDeleteLogic {
	return &FileDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FileDeleteLogic) FileDelete(in *pb.FileIdRequest) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("id 不能为空")
	}
	// 逻辑删除（deleted=1），删除不存在的 id 幂等返回成功
	if err := l.svcCtx.FileModel.SoftDelete(l.ctx, in.Id); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "删除文件失败")
	}
	return &pb.Empty{}, nil
}
