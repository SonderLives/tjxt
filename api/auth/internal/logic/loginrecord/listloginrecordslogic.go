// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package loginrecord

import (
	"context"

	backendlogic "tjxt/api/auth/backend/logic"
	authpb "tjxt/api/auth/backend/pb"
	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoginRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLoginRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginRecordsLogic {
	return &ListLoginRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListLoginRecords 分页查询登录记录，可按 userId 过滤（0 表示查全部）。
func (l *ListLoginRecordsLogic) ListLoginRecords(req *types.LoginRecordListReq) (resp *types.LoginRecordListVO, err error) {
	reply, err := backendlogic.NewListLoginRecordsLogic(l.ctx, l.svcCtx.Backend).ListLoginRecords(&authpb.LoginRecordPageReq{
		PageNo:   int32(req.PageNo),
		PageSize: int32(req.PageSize),
		UserId:   req.UserId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.LoginRecordVO, 0, len(reply.List))
	for _, v := range reply.List {
		list = append(list, types.LoginRecordVO{
			Id:         v.Id,
			UserId:     v.UserId,
			CellPhone:  v.CellPhone,
			LoginTime:  v.LoginTime,
			LogoutTime: v.LogoutTime,
			Duration:   v.Duration,
			Ipv4:       v.Ipv4,
		})
	}
	return &types.LoginRecordListVO{Total: reply.Total, List: list}, nil
}
