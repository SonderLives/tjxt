package datalogic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"tjxt/service/data/api/backend/pb"
	"tjxt/service/data/api/backend/svc"
	"tjxt/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

// boardKeyFormat 看板数据 Redis key 模板（按 type 区分）
const boardKeyFormat = "data:board:%d"

// versionKeySuffix 版本键后缀（version 预留，用于未来多版本原子切换）
const versionKeySuffix = ":version"

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
		return nil, errorx.BadRequestf("看板类型不合法")
	}
	if len(in.Data) == 0 {
		return nil, errorx.BadRequestf("看板数据不能为空")
	}

	data, err := json.Marshal(in.Data)
	if err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "序列化看板数据失败")
	}
	key := fmt.Sprintf(boardKeyFormat, in.Type)
	if err := l.svcCtx.Rds.SetCtx(l.ctx, key, string(data)); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "写入看板数据失败")
	}
	// version 预留,用于未来多版本原子切换;读侧不参与逻辑,仅旁路落盘保留
	if err := l.svcCtx.Rds.SetCtx(l.ctx, key+versionKeySuffix, strconv.Itoa(int(in.Version))); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "写入看板版本失败")
	}
	return &pb.OkReply{Success: true}, nil
}
