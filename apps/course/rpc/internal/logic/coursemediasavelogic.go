package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseMediaSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseMediaSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseMediaSaveLogic {
	return &CourseMediaSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 课程媒体 =====
// CourseMediaSave 把媒资信息回填到课程草稿目录（小节）上。
func (l *CourseMediaSaveLogic) CourseMediaSave(in *pb.CourseMediaSaveRequest) (*pb.Empty, error) {
	if in.CourseId == 0 {
		return nil, xerr.BadRequestf("课程id不能为空")
	}
	if len(in.Medias) == 0 {
		return &pb.Empty{}, nil
	}

	catalogues, err := l.svcCtx.CourseCatalogueDraftModel.ListByCourseId(l.ctx, in.CourseId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程目录草稿失败")
	}
	cataMap := make(map[int64]*model.CourseCatalogueDraft, len(catalogues))
	for _, c := range catalogues {
		cataMap[c.Id] = c
	}

	for _, media := range in.Medias {
		if media == nil {
			continue
		}
		cata, ok := cataMap[media.CataId]
		if !ok {
			return nil, xerr.NotFound("课程目录不存在")
		}
		cata.MediaId = media.MediaId
		cata.VideoName = media.VideoName
		cata.MediaDuration = media.MediaDuration
		if media.Trailer {
			cata.Trailer = 1
		} else {
			cata.Trailer = 0
		}
		if err := l.svcCtx.CourseCatalogueDraftModel.Update(l.ctx, cata); err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "保存课程媒资失败")
		}
	}
	return &pb.Empty{}, nil
}
