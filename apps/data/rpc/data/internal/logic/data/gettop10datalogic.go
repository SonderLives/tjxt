package datalogic

import (
	"context"
	"encoding/json"
	"sort"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const top10Limit = 10

type GetTop10DataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTop10DataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTop10DataLogic {
	return &GetTop10DataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTop10DataLogic) GetTop10Data(in *pb.Empty) (*pb.Top10DataVO, error) {
	val, err := l.svcCtx.Rds.GetCtx(l.ctx, top10DataKey)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "读取榜单数据失败")
	}
	// key 不存在返回空结构
	if val == "" {
		return &pb.Top10DataVO{}, nil
	}

	var units []*pb.Top10DataSetUnit
	if err := json.Unmarshal([]byte(val), &units); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "解析榜单数据失败")
	}

	// 同一份数据按不同指标派生出两个榜单:hot 按新增学员数,hotSales 按订单额
	hot := append([]*pb.Top10DataSetUnit(nil), units...)
	sort.SliceStable(hot, func(i, j int) bool {
		return hot[i].NewStuNum > hot[j].NewStuNum
	})
	hotSales := append([]*pb.Top10DataSetUnit(nil), units...)
	sort.SliceStable(hotSales, func(i, j int) bool {
		return hotSales[i].OrderAmount > hotSales[j].OrderAmount
	})

	resp := &pb.Top10DataVO{}
	for _, u := range limitTop10(hot) {
		resp.Hot = append(resp.Hot, toCourseInfo(u))
	}
	for _, u := range limitTop10(hotSales) {
		resp.HotSales = append(resp.HotSales, toCourseInfo(u))
	}
	return resp, nil
}

func limitTop10(units []*pb.Top10DataSetUnit) []*pb.Top10DataSetUnit {
	if len(units) > top10Limit {
		return units[:top10Limit]
	}
	return units
}

func toCourseInfo(u *pb.Top10DataSetUnit) *pb.CourseInfo {
	return &pb.CourseInfo{
		Category:    u.Category,
		Name:        u.Name,
		NewStuNum:   u.NewStuNum,
		OrderAmount: u.OrderAmount,
	}
}
