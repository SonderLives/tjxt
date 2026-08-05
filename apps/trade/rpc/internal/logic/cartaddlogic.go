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

type CartAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartAddLogic {
	return &CartAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 购物车 =====
func (l *CartAddLogic) CartAdd(in *pb.CartAddRequest) (*pb.Empty, error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}
	if in.CourseId <= 0 {
		return nil, xerr.BadRequestf("课程ID不能为空")
	}

	// 去重：同一用户不能重复加入同一课程
	if _, e := l.svcCtx.CartModel.FindByUserIdAndCourseId(l.ctx, userId, in.CourseId); e == nil {
		return nil, xerr.New(xerr.CodeConflict, "该课程已在购物车中")
	} else if !errors.Is(e, model.ErrNotFound) {
		return nil, xerr.Wrap(e, xerr.CodeInternal, "查询购物车失败")
	}

	// 校验课程存在并落库名称/封面/价格
	courseMap := fetchCourseMap(l.ctx, l.svcCtx, []int64{in.CourseId})
	c, ok := courseMap[in.CourseId]
	if !ok {
		return nil, xerr.New(xerr.CodeNotFound, "课程不存在")
	}

	if _, err = l.svcCtx.CartModel.Insert(l.ctx, &model.Cart{
		Id:         nextID(),
		UserId:     userId,
		CourseId:   in.CourseId,
		CourseName: c.Name,
		CoverUrl:   c.CoverUrl,
		Price:      c.Price,
		CreateTime: now(),
	}); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "加入购物车失败")
	}
	return &pb.Empty{}, nil
}
