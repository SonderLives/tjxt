// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/api/course/internal/svc"
	"tjxt/api/course/internal/types"
	"tjxt/rpc/course/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CatalogueSectionInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCatalogueSectionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CatalogueSectionInfoLogic {
	return &CatalogueSectionInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CatalogueSectionInfo 查询小节信息（所属课程与媒资）。
func (l *CatalogueSectionInfoLogic) CatalogueSectionInfo(req *types.IdPathReq) (resp *types.CourseSectionInfoVO, err error) {
	info, gerr := l.svcCtx.CourseRpc.CourseCatalogueSectionInfo(l.ctx, &pb.IdRequest{Id: req.Id})
	if gerr != nil {
		return nil, gerr
	}
	return &types.CourseSectionInfoVO{
		Id:       info.Id,
		CourseId: info.CourseId,
		MediaId:  info.MediaId,
	}, nil
}
