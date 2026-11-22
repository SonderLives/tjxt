package userlogic

import (
	"context"

	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryStaffsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageQueryStaffsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryStaffsLogic {
	return &PageQueryStaffsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PageQueryStaffs 员工分页查询（type=1）。
func (l *PageQueryStaffsLogic) PageQueryStaffs(in *pb.UserPageRequest) (*pb.StaffPageResponse, error) {
	offset, limit := page.Normalize(in.PageNo, in.PageSize)
	col, dir := sortClause(in.SortBy, in.IsAsc)
	list, total, err := l.svcCtx.UserModel.FindPageByType(l.ctx, 1, in.Name, in.Phone, int64(in.Status), offset, limit, col, dir)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询员工列表失败")
	}
	items := make([]*pb.StaffVO, 0, len(list))
	for _, v := range list {
		items = append(items, toStaffVO(v))
	}
	return &pb.StaffPageResponse{
		Total: total,
		Pages: page.CalcPages(total, limit),
		List:  items,
	}, nil
}
