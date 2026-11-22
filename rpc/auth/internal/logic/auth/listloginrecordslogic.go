package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/utils/page"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoginRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLoginRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginRecordsLogic {
	return &ListLoginRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListLoginRecords 分页查询登录记录，UserId 为 0 时查询全部用户。
func (l *ListLoginRecordsLogic) ListLoginRecords(in *pb.LoginRecordPageReq) (*pb.LoginRecordListReply, error) {
	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))

	records, total, err := l.svcCtx.LoginRecordModel.FindPage(l.ctx, in.UserId, offset, limit)
	if err != nil {
		return nil, err
	}

	list := make([]*pb.LoginRecordVO, 0, len(records))
	for _, r := range records {
		list = append(list, toLoginRecordVO(r))
	}
	return &pb.LoginRecordListReply{Total: total, List: list}, nil
}
