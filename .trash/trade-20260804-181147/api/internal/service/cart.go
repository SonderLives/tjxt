package service

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/xerr"
	"tjxt/apps/trade/api/internal/model"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// CartService 购物车业务接口
type CartService interface {
	// AddToCart 添加课程到购物车（幂等）。
	AddToCart(ctx context.Context, userId, courseId int64) error
	// ListCart 查询用户购物车。
	ListCart(ctx context.Context, userId int64) ([]types.CartVO, error)
	// DeleteCart 删除购物车条目。
	DeleteCart(ctx context.Context, userId int64, ids []int64) error
}

type cartService struct {
	cartModel    *model.CartModel
	courseClient CourseClient
}

// NewCartService 创建购物车业务服务。
func NewCartService(cartModel *model.CartModel, courseClient CourseClient) CartService {
	return &cartService{cartModel: cartModel, courseClient: courseClient}
}

func (s *cartService) AddToCart(ctx context.Context, userId, courseId int64) error {
	if userId == 0 || courseId == 0 {
		return xerr.BadRequestf("课程 id 不能为空")
	}

	// 校验课程存在
	infos, err := s.courseClient.GetSimpleInfos(ctx, []int64{courseId})
	if err != nil {
		return err
	}
	info, ok := infos[courseId]
	if !ok || info == nil {
		return xerr.NotFound("课程不存在")
	}

	// 幂等：已在购物车则直接返回
	existing, err := s.cartModel.FindByUserCourse(ctx, userId, courseId)
	if err != nil && err != sql.ErrNoRows {
		return xerr.Wrap(err, xerr.CodeInternal, "查询购物车失败")
	}
	if existing != nil {
		return nil
	}

	item := &model.Cart{
		Id:         idgen.NextID(),
		UserId:     userId,
		CourseId:   courseId,
		CoverUrl:   info.CoverUrl,
		CourseName: info.Name,
		Price:      info.Price,
	}
	if err := s.cartModel.Insert(ctx, item); err != nil {
		logx.Errorf("add cart failed, user=%d course=%d err=%v", userId, courseId, err)
		return xerr.Wrap(err, xerr.CodeInternal, "加入购物车失败")
	}
	return nil
}

func (s *cartService) ListCart(ctx context.Context, userId int64) ([]types.CartVO, error) {
	rows, err := s.cartModel.ListByUser(ctx, userId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询购物车失败")
	}
	if len(rows) == 0 {
		return []types.CartVO{}, nil
	}

	courseIds := make([]int64, 0, len(rows))
	for i := range rows {
		courseIds = append(courseIds, rows[i].CourseId)
	}
	infos, err := s.courseClient.GetSimpleInfos(ctx, courseIds)
	if err != nil {
		return nil, err
	}

	list := make([]types.CartVO, 0, len(rows))
	for i := range rows {
		row := &rows[i]
		vo := types.CartVO{
			Id:         row.Id,
			CourseId:   row.CourseId,
			CourseName: row.CourseName,
			CoverUrl:   row.CoverUrl,
			Price:      row.Price,
			NowPrice:   row.Price,
			Expired:    true,
		}
		if info, ok := infos[row.CourseId]; ok && info != nil {
			vo.NowPrice = info.Price
			// 课程状态为下架(2)时标记过期
			vo.Expired = info.Status == 2
		}
		list = append(list, vo)
	}
	return list, nil
}

func (s *cartService) DeleteCart(ctx context.Context, userId int64, ids []int64) error {
	if len(ids) == 0 {
		return xerr.BadRequestf("请选择要删除的课程")
	}
	if err := s.cartModel.DeleteByIds(ctx, userId, ids); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除购物车失败")
	}
	return nil
}

// ParseIDs 将逗号分隔的字符串解析为 id 列表。
func ParseIDs(raw string) ([]int64, error) {
	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, xerr.BadRequestf("非法的 id: %s", p)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
