package datalogic

import (
	"context"
	"encoding/json"
	"strconv"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const top10DataKey = "data:top10"

type SetTop10DataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetTop10DataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetTop10DataLogic {
	return &SetTop10DataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetTop10DataLogic) SetTop10Data(in *pb.Top10DataSetReq) (*pb.OkReply, error) {
	if len(in.Data) == 0 {
		return nil, xerr.BadRequestf("榜单数据不能为空")
	}

	data, err := json.Marshal(in.Data)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "序列化榜单数据失败")
	}
	if err := l.svcCtx.Rds.SetCtx(l.ctx, top10DataKey, string(data)); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "写入榜单数据失败")
	}
	// version 预留,用于未来多版本原子切换;读侧不参与逻辑,仅旁路落盘保留
	if err := l.svcCtx.Rds.SetCtx(l.ctx, top10DataKey+":version", strconv.Itoa(int(in.Version))); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "写入榜单版本失败")
	}
	return &pb.OkReply{Success: true}, nil
}
