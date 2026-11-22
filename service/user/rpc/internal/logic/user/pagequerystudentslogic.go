package userlogic

import (
	"context"

	"tjxt/service/user/rpc/internal/svc"
	"tjxt/service/user/rpc/pb"
	"tjxt/common/page"
	"tjxt/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryStudentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageQueryStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryStudentsLogic {
	return &PageQueryStudentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PageQueryStudents 学员分页查询（type=2）。
func (l *PageQueryStudentsLogic) PageQueryStudents(in *pb.UserPageRequest) (*pb.StudentPageResponse, error) {
	offset, limit := page.Normalize(in.PageNo, in.PageSize)
	col, dir := sortClause(in.SortBy, in.IsAsc)
	list, total, err := l.svcCtx.UserModel.FindPageByType(l.ctx, 2, in.Name, in.Phone, int64(in.Status), offset, limit, col, dir)
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询学员列表失败")
	}
	items := make([]*pb.StudentPageVO, 0, len(list))
	for _, v := range list {
		items = append(items, toStudentPageVO(v))
	}
	return &pb.StudentPageResponse{
		Total: total,
		Pages: page.CalcPages(total, limit),
		List:  items,
	}, nil
}
