package datalogic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBoardDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBoardDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBoardDataLogic {
	return &GetBoardDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBoardDataLogic) GetBoardData(in *pb.BoardDataReq) (*pb.EchartsVO, error) {
	if len(in.Types) == 0 {
		return &pb.EchartsVO{}, nil
	}

	resp := &pb.EchartsVO{}
	maxLen := 0
	for _, t := range in.Types {
		vals, err := l.readBoardData(int(t))
		if err != nil {
			return nil, err
		}
		if len(vals) > maxLen {
			maxLen = len(vals)
		}
		resp.Series = append(resp.Series, buildSerier(t, vals))
		resp.YAxis = append(resp.YAxis, buildYAxis(vals))
	}
	resp.XAxis = []*pb.AxisVO{buildXAxis(maxLen)}
	return resp, nil
}

func (l *GetBoardDataLogic) readBoardData(t int) ([]float64, error) {
	val, err := l.svcCtx.Rds.GetCtx(l.ctx, fmt.Sprintf("data:board:%d", t))
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "读取看板数据失败")
	}
	if val == "" {
		return nil, nil
	}
	var vals []float64
	if err := json.Unmarshal([]byte(val), &vals); err != nil {
		l.Logger.Errorf("解析看板数据失败 type=%d err=%v", t, err)
		return nil, nil
	}
	return vals, nil
}

// buildSerier 每个 type 一个序列:名称即 type,数据逐项转字符串,max/min 为该序列极值。
func buildSerier(t int32, vals []float64) *pb.SerierVO {
	s := &pb.SerierVO{Name: fmt.Sprintf("%d", t), Type: "line"}
	for _, v := range vals {
		s.Data = append(s.Data, strconv.FormatFloat(v, 'f', -1, 64))
	}
	if len(vals) > 0 {
		max, min := vals[0], vals[0]
		for _, v := range vals[1:] {
			if v > max {
				max = v
			}
			if v < min {
				min = v
			}
		}
		s.Max = strconv.FormatFloat(max, 'f', -1, 64)
		s.Min = strconv.FormatFloat(min, 'f', -1, 64)
	}
	return s
}

// buildXAxis 分类轴,索引字符串序列,长度对齐最长序列。
func buildXAxis(maxLen int) *pb.AxisVO {
	axis := &pb.AxisVO{Type: "category"}
	for i := 0; i < maxLen; i++ {
		axis.Data = append(axis.Data, strconv.Itoa(i))
	}
	return axis
}

// buildYAxis 数值轴,为该序列现算 max/min/average/interval。
func buildYAxis(vals []float64) *pb.AxisVO {
	axis := &pb.AxisVO{Type: "value"}
	if len(vals) == 0 {
		return axis
	}
	max, min, sum := vals[0], vals[0], vals[0]
	for _, v := range vals[1:] {
		sum += v
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	axis.Max = max
	axis.Min = min
	axis.Average = sum / float64(len(vals))
	if len(vals) >= 2 {
		axis.Interval = vals[1] - vals[0]
	}
	return axis
}
