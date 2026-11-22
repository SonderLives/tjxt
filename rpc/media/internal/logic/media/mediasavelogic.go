package medialogic

import (
	"context"

	"tjxt/rpc/media/internal/model"
	"tjxt/rpc/media/internal/svc"
	"tjxt/rpc/media/pb"
	"tjxt/pkg/idgen"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMediaSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaSaveLogic {
	return &MediaSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MediaSaveLogic) MediaSave(in *pb.MediaSaveRequest) (*pb.MediaIdReply, error) {
	if in.Id > 0 {
		return l.update(in)
	}
	return l.create(in)
}

// create 新建媒资：校验文件名与关联文件（status=2），创建后把文件标记为已使用
func (l *MediaSaveLogic) create(in *pb.MediaSaveRequest) (*pb.MediaIdReply, error) {
	if in.Filename == "" {
		return nil, errorx.BadRequestf("文件名不能为空")
	}
	if in.FileId == "" {
		return nil, errorx.BadRequestf("fileId 不能为空")
	}
	file, err := l.svcCtx.FileModel.FindByKey(l.ctx, in.FileId)
	if err != nil {
		if isNotFound(err) {
			return nil, errorx.NotFound("文件不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询文件失败")
	}
	if file.Status != FileStatusUploaded {
		return nil, errorx.Conflict("文件尚未上传或已被使用")
	}

	id := idgen.NextID()
	// 配置对象存储时生成公共读持久地址，否则 mock 占位
	mediaURL := mockFileURL(in.FileId)
	if l.svcCtx.Store != nil {
		mediaURL = l.svcCtx.Store.PublicURL(in.FileId)
	}
	media := &model.Media{
		Id:       id,
		FileId:   in.FileId,
		Filename: in.Filename,
		MediaUrl: mediaURL,
		Duration: in.Duration,
		Size:     in.Size,
		Status:   MediaStatusUploaded,
		// RPC 协议暂无用户字段，creater/updater 保持默认 0
	}
	if _, err := l.svcCtx.MediaModel.Insert(l.ctx, media); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "创建媒资失败")
	}
	// 媒资关联成功后把文件标记为已使用（status=3），失败仅记录日志不影响主流程
	if err := l.svcCtx.FileModel.UpdateStatus(l.ctx, file.Id, FileStatusUsed); err != nil {
		logx.WithContext(l.ctx).Errorf("更新文件状态失败: fileId=%d err=%v", file.Id, err)
	}
	return &pb.MediaIdReply{Id: id}, nil
}

// update 更新媒资：只允许更新 filename/duration/size，file_id 不允许变更
func (l *MediaSaveLogic) update(in *pb.MediaSaveRequest) (*pb.MediaIdReply, error) {
	media, err := l.svcCtx.MediaModel.FindOneNotDeleted(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, errorx.NotFound("媒资不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询媒资失败")
	}
	if in.Filename != "" {
		media.Filename = in.Filename
	}
	media.Duration = in.Duration
	media.Size = in.Size
	// RPC 协议无 media_url/cover_url 字段，保持原值不变
	if err := l.svcCtx.MediaModel.Update(l.ctx, media); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "更新媒资失败")
	}
	return &pb.MediaIdReply{Id: media.Id}, nil
}
