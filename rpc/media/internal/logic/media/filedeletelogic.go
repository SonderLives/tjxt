package medialogic

import (
	"context"

	"tjxt/rpc/media/internal/svc"
	"tjxt/rpc/media/pb"
	"tjxt/pkg/errorx"

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
		return nil, errorx.BadRequestf("id 不能为空")
	}
	// 逻辑删除（deleted=1），删除不存在的 id 幂等返回成功
	if err := l.svcCtx.FileModel.SoftDelete(l.ctx, in.Id); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "删除文件失败")
	}
	return &pb.Empty{}, nil
}
