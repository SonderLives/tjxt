package logic

import (
	"context"

	"tjxt/service/course/api/internal/svc"
	"tjxt/service/course/api/internal/types"
	"tjxt/service/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryGetLogic {
	return &CategoryGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CategoryGet 查询分类详情。
func (l *CategoryGetLogic) CategoryGet(req *types.IdPathReq) (resp *types.CategoryInfoVO, err error) {
	info, gerr := l.svcCtx.CourseRpc.CategoryGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if gerr != nil {
		return nil, gerr
	}
	return &types.CategoryInfoVO{
		Id:                 info.Id,
		Name:               info.Name,
		Level:              int64(info.Level),
		FirstCategoryName:  info.FirstCategoryName,
		SecondCategoryName: info.SecondCategoryName,
		Status:             int64(info.Status),
		Index:              int64(info.Priority),
		CreateTime:         info.CreateTime,
		UpdateTime:         info.UpdateTime,
	}, nil
}
