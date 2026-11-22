// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/service/course/api/internal/svc"
	"tjxt/service/course/api/internal/types"
	"tjxt/service/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseMediaSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseMediaSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseMediaSaveLogic {
	return &CourseMediaSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseMediaSave 保存课程目录媒资（透传 RPC）。
func (l *CourseMediaSaveLogic) CourseMediaSave(req *types.CourseMediaSaveReq) (resp *types.NameExistVO, err error) {
	medias := make([]*pb.CourseMediaBind, 0, len(req.List))
	for _, item := range req.List {
		medias = append(medias, &pb.CourseMediaBind{
			CataId:        item.CataId,
			MediaId:       item.MediaId,
			VideoName:     item.VideoName,
			MediaDuration: item.MediaDuration,
			Trailer:       item.Trailer,
		})
	}
	_, err = l.svcCtx.CourseRpc.CourseMediaSave(l.ctx, &pb.CourseMediaSaveRequest{
		CourseId: req.Id,
		Medias:   medias,
	})
	if err != nil {
		return nil, err
	}
	return &types.NameExistVO{Existed: false}, nil
}
