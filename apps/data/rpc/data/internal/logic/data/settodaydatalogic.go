package datalogic

import (
	"context"
	"strconv"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const todayDataKey = "data:today"

type SetTodayDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetTodayDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetTodayDataLogic {
	return &SetTodayDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetTodayDataLogic) SetTodayData(in *pb.TodayDataSetReq) (*pb.OkReply, error) {
	// version 预留,用于未来多版本原子切换;读侧不参与逻辑,仅作为字段落盘保留
	fields := map[string]string{
		"visits":      strconv.FormatFloat(in.Visits, 'f', -1, 64),
		"orderAmount": strconv.FormatFloat(in.OrderAmount, 'f', -1, 64),
		"orderNum":    strconv.Itoa(int(in.OrderNum)),
		"stuNewNum":   strconv.Itoa(int(in.StuNewNum)),
		"version":     strconv.Itoa(int(in.Version)),
	}
	if err := l.svcCtx.Rds.HmsetCtx(l.ctx, todayDataKey, fields); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "写入今日数据失败")
	}
	return &pb.OkReply{Success: true}, nil
}
