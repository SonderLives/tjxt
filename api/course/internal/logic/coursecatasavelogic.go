// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCataSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCataSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCataSaveLogic {
	return &CourseCataSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseCataSave 保存课程目录（章节树全量覆盖）。
func (l *CourseCataSaveLogic) CourseCataSave(req *types.CataSaveReq) (resp *types.NameExistVO, err error) {
	chapters := make([]*pb.CourseChapterSave, 0, len(req.List))
	for i := range req.List {
		chapters = append(chapters, cataSaveToChapterSave(&req.List[i]))
	}
	_, err = l.svcCtx.CourseRpc.CourseCatalogueSave(l.ctx, &pb.CourseCatalogueSaveRequest{
		CourseId: req.Id,
		Step:     req.Step,
		Chapters: chapters,
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "保存课程目录失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}

// cataSaveToChapterSave API CataSaveItemReq -> pb.CourseChapterSave（递归映射子节）。
func cataSaveToChapterSave(item *types.CataSaveItemReq) *pb.CourseChapterSave {
	if item == nil {
		return nil
	}
	save := &pb.CourseChapterSave{
		Id:    item.Id,
		Name:  item.Name,
		Index: int32(item.Index),
		Type:  int32(item.Type),
	}
	for i := range item.Sections {
		save.Sections = append(save.Sections, cataSaveToChapterSave(&item.Sections[i]))
	}
	return save
}
