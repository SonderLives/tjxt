// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseMediaUseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseMediaUseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseMediaUseInfoLogic {
	return &CourseMediaUseInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseMediaUseInfo 查询媒资被课程引用的数量（透传 RPC）。
func (l *CourseMediaUseInfoLogic) CourseMediaUseInfo(req *types.CourseMediaUseInfoReq) (resp []types.MediaQuoteVO, err error) {
	list, gerr := l.svcCtx.CourseRpc.CourseMediaUseInfo(l.ctx, &pb.MediaIdsRequest{
		MediaIds: parseIds(req.MediaIds),
	})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询媒资引用信息失败")
	}
	resp = make([]types.MediaQuoteVO, 0, len(list.Items))
	for _, item := range list.Items {
		resp = append(resp, types.MediaQuoteVO{
			MediaId:  item.MediaId,
			QuoteNum: item.QuoteNum,
		})
	}
	return resp, nil
}
