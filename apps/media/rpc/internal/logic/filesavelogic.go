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

type FileSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileSaveLogic {
	return &FileSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FileSaveLogic) FileSave(in *pb.FileSaveRequest) (*pb.FileIdReply, error) {
	if in.Filename == "" || in.Key == "" {
		return nil, xerr.BadRequestf("filename/key 不能为空")
	}
	// 按 key 查重，已存在（未删除）返回冲突
	if _, err := l.svcCtx.FileModel.FindByKey(l.ctx, in.Key); err == nil {
		return nil, xerr.Conflict("文件已存在")
	} else if !isNotFound(err) {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询文件失败")
	}

	id := idgen.NextID()
	file := &model.File{
		Id:       id,
		Key:      in.Key,
		Filename: in.Filename,
		Status:   FileStatusPending,
		Platform: PlatformTencent,
	}
	if _, err := l.svcCtx.FileModel.Insert(l.ctx, file); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "创建文件记录失败")
	}
	return &pb.FileIdReply{Id: id}, nil
}
