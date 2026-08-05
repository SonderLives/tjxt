package logic

import (
	"context"
	"errors"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartUpdateLogic {
	return &CartUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CartUpdateLogic) CartUpdate(in *pb.CartUpdateRequest) (*pb.Empty, error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("购物车条目ID不能为空")
	}
	if in.CourseId <= 0 {
		return nil, xerr.BadRequestf("课程ID不能为空")
	}

	c, err := l.svcCtx.CartModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("购物车条目不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询购物车条目失败")
	}
	if c.UserId != userId {
		return nil, xerr.NotFound("购物车条目不存在")
	}

	// 课程变更后同步刷新课程名称/封面/价格快照
	if course, ok := fetchCourseMap(l.ctx, l.svcCtx, []int64{in.CourseId})[in.CourseId]; ok {
		c.CourseName = course.Name
		c.CoverUrl = course.CoverUrl
		c.Price = course.Price
	}
	c.CourseId = in.CourseId

	if err = l.svcCtx.CartModel.Update(l.ctx, c); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新购物车条目失败")
	}
	return &pb.Empty{}, nil
}
