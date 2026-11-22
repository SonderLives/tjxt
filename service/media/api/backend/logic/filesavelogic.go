package medialogic

import (
	"context"
	"errors"

	"tjxt/service/media/api/backend/model"
	"tjxt/service/media/api/backend/pb"
	"tjxt/service/media/api/backend/svc"
	"tjxt/common/errorx"
	"tjxt/common/idgen"

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
		return nil, errorx.BadRequestf("filename/key 不能为空")
	}
	// 按 key 查重，已存在（未删除）返回冲突
	if _, err := l.svcCtx.FileModel.FindByKey(l.ctx, in.Key); err == nil {
		return nil, errorx.Conflict("文件已存在")
	} else if !errors.Is(err, model.ErrNotFound) {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询文件失败")
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
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "创建文件记录失败")
	}
	return &pb.FileIdReply{Id: id}, nil
}
