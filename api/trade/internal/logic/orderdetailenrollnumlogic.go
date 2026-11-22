// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"strconv"
	"strings"

	backendlogic "tjxt/api/trade/backend/logic"
	"tjxt/api/trade/backend/pb"
	"tjxt/api/trade/internal/svc"
	"tjxt/api/trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailEnrollNumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailEnrollNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailEnrollNumLogic {
	return &OrderDetailEnrollNumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDetailEnrollNumLogic) OrderDetailEnrollNum(req *types.EnrollNumReq) (resp *types.NamePlaceVO, err error) {
	var ids []int64
	for _, s := range strings.Split(req.CourseIdList, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, parseErr := strconv.ParseInt(s, 10, 64)
		if parseErr != nil {
			continue
		}
		ids = append(ids, id)
	}

	if _, err = backendlogic.NewOrderDetailEnrollNumLogic(l.ctx, l.svcCtx.Backend).OrderDetailEnrollNum(&pb.EnrollNumRequest{CourseIdList: ids}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
