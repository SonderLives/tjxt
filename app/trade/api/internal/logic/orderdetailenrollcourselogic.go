// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"strconv"
	"strings"

	backendlogic "tjxt/service/trade/api/backend/logic"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/service/trade/api/internal/svc"
	"tjxt/service/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailEnrollCourseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailEnrollCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailEnrollCourseLogic {
	return &OrderDetailEnrollCourseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDetailEnrollCourseLogic) OrderDetailEnrollCourse(req *types.EnrollCourseReq) (resp *types.NamePlaceVO, err error) {
	var ids []int64
	for _, s := range strings.Split(req.StudentIds, ",") {
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

	reply, err := backendlogic.NewOrderDetailEnrollCourseLogic(l.ctx, l.svcCtx.Backend).OrderDetailEnrollCourse(&pb.EnrollCourseRequest{StudentIds: ids})
	if err != nil {
		return nil, err
	}

	return &types.NamePlaceVO{
		Existed: true,
		Message: "ok:" + strconv.Itoa(len(reply.Items)),
	}, nil
}
