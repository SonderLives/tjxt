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

type CourseBaseInfoSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseBaseInfoSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseBaseInfoSaveLogic {
	return &CourseBaseInfoSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseBaseInfoSave 保存课程基础信息（透传 RPC）。
func (l *CourseBaseInfoSaveLogic) CourseBaseInfoSave(req *types.CourseBaseInfoSaveReq) (resp *types.CourseSaveVO, err error) {
	res, gerr := l.svcCtx.CourseRpc.CourseBaseInfoSave(l.ctx, &pb.CourseBaseInfoSaveRequest{
		Id:                req.Id,
		Name:              req.Name,
		CoverUrl:          req.CoverUrl,
		Price:             req.Price,
		Free:              int32(req.Free),
		ThirdCateId:       req.ThirdCateId,
		Introduce:         req.Introduce,
		Detail:            req.Detail,
		UsePeople:         req.UsePeople,
		ValidDuration:     req.ValidDuration,
		PurchaseStartTime: req.PurchaseStartTime,
		PurchaseEndTime:   req.PurchaseEndTime,
	})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "保存课程基础信息失败")
	}
	return &types.CourseSaveVO{Id: res.Id}, nil
}
