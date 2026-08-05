package datalogic

import (
	"context"
	"strconv"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTodayDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTodayDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTodayDataLogic {
	return &GetTodayDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTodayDataLogic) GetTodayData(in *pb.Empty) (*pb.TodayDataVO, error) {
	vals, err := l.svcCtx.Rds.HgetallCtx(l.ctx, todayDataKey)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "读取今日数据失败")
	}
	// key 不存在返回全零结构,大屏首启友好
	resp := &pb.TodayDataVO{}
	if len(vals) == 0 {
		return resp, nil
	}
	if v, err := strconv.ParseFloat(vals["visits"], 64); err == nil {
		resp.Visits = v
	} else {
		l.Logger.Errorf("解析今日数据失败 field=visits value=%q err=%v", vals["visits"], err)
	}
	if v, err := strconv.ParseFloat(vals["orderAmount"], 64); err == nil {
		resp.OrderAmount = v
	} else {
		l.Logger.Errorf("解析今日数据失败 field=orderAmount value=%q err=%v", vals["orderAmount"], err)
	}
	if v, err := strconv.ParseInt(vals["orderNum"], 10, 32); err == nil {
		resp.OrderNum = int32(v)
	} else {
		l.Logger.Errorf("解析今日数据失败 field=orderNum value=%q err=%v", vals["orderNum"], err)
	}
	if v, err := strconv.ParseInt(vals["stuNewNum"], 10, 32); err == nil {
		resp.StuNewNum = int32(v)
	} else {
		l.Logger.Errorf("解析今日数据失败 field=stuNewNum value=%q err=%v", vals["stuNewNum"], err)
	}
	return resp, nil
}
