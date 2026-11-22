package tradelogic

import (
	"context"
	"errors"

	"tjxt/rpc/trade/internal/model"
	"tjxt/rpc/trade/internal/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/errorx"

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
		return nil, errorx.New(errorx.CodeUnauthorized, "未登录")
	}
	if in.Id <= 0 {
		return nil, errorx.BadRequestf("购物车条目ID不能为空")
	}
	if in.CourseId <= 0 {
		return nil, errorx.BadRequestf("课程ID不能为空")
	}

	c, err := l.svcCtx.CartModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("购物车条目不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询购物车条目失败")
	}
	if c.UserId != userId {
		return nil, errorx.NotFound("购物车条目不存在")
	}

	// 课程变更后同步刷新课程名称/封面/价格快照
	if course, ok := fetchCourseMap(l.ctx, l.svcCtx, []int64{in.CourseId})[in.CourseId]; ok {
		c.CourseName = course.Name
		c.CoverUrl = course.CoverUrl
		c.Price = course.Price
	}
	c.CourseId = in.CourseId

	if err = l.svcCtx.CartModel.Update(l.ctx, c); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "更新购物车条目失败")
	}
	return &pb.Empty{}, nil
}
