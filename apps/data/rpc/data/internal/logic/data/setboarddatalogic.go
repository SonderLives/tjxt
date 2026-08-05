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

type SetBoardDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetBoardDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetBoardDataLogic {
	return &SetBoardDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetBoardDataLogic) SetBoardData(in *pb.BoardDataSetReq) (*pb.OkReply, error) {
	if in.Type <= 0 {
		return nil, xerr.BadRequestf("看板类型不合法")
	}
	if len(in.Data) == 0 {
		return nil, xerr.BadRequestf("看板数据不能为空")
	}

	data, err := json.Marshal(in.Data)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "序列化看板数据失败")
	}
	key := fmt.Sprintf("data:board:%d", in.Type)
	if err := l.svcCtx.Rds.SetCtx(l.ctx, key, string(data)); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "写入看板数据失败")
	}
	// version 预留,用于未来多版本原子切换;读侧不参与逻辑,仅旁路落盘保留
	if err := l.svcCtx.Rds.SetCtx(l.ctx, key+":version", strconv.Itoa(int(in.Version))); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "写入看板版本失败")
	}
	return &pb.OkReply{Success: true}, nil
}
