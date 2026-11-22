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

type CourseBaseInfoGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseBaseInfoGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseBaseInfoGetLogic {
	return &CourseBaseInfoGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseBaseInfoGet 查询课程基础信息（透传 RPC）。
func (l *CourseBaseInfoGetLogic) CourseBaseInfoGet(req *types.CourseBaseInfoGetReq) (resp *types.CourseBaseInfoVO, err error) {
	info, gerr := l.svcCtx.CourseRpc.CourseBaseInfoGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询课程基础信息失败")
	}
	return &types.CourseBaseInfoVO{
		Id:                info.Id,
		Name:              info.Name,
		Price:             info.Price,
		Free:              int64(info.Free),
		CoverUrl:          info.CoverUrl,
		Detail:            info.Detail,
		Introduce:         info.Introduce,
		UsePeople:         info.UsePeople,
		ValidDuration:     info.ValidDuration,
		PurchaseStartTime: info.PurchaseStartTime,
		PurchaseEndTime:   info.PurchaseEndTime,
		FirstCateId:       info.FirstCateId,
		SecondCateId:      info.SecondCateId,
		ThirdCateId:       info.ThirdCateId,
		CateNames:         info.CateNames,
		Status:            int64(info.Status),
		Step:              int64(info.Step),
		Score:             info.Score,
		CourseScore:       info.CourseScore,
		CataTotalNum:      info.CataTotalNum,
		EnrollNum:         info.EnrollNum,
		StudyNum:          info.StudyNum,
		RealPayAmount:     info.RealPayAmount,
		RefundNum:         info.RefundNum,
		CreaterName:       info.CreaterName,
		UpdaterName:       info.UpdaterName,
		CanUpdate:         info.CanUpdate,
		CreateTime:        info.CreateTime,
		UpdateTime:        info.UpdateTime,
	}, nil
}
