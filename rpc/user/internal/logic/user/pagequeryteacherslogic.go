package userlogic

import (
	"context"

	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryTeachersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageQueryTeachersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryTeachersLogic {
	return &PageQueryTeachersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PageQueryTeachers 老师分页查询（type=3）。
func (l *PageQueryTeachersLogic) PageQueryTeachers(in *pb.UserPageRequest) (*pb.TeacherPageResponse, error) {
	offset, limit := page.Normalize(in.PageNo, in.PageSize)
	col, dir := sortClause(in.SortBy, in.IsAsc)
	list, total, err := l.svcCtx.UserModel.FindPageByType(l.ctx, 3, in.Name, in.Phone, int64(in.Status), offset, limit, col, dir)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询老师列表失败")
	}
	items := make([]*pb.TeacherPageVO, 0, len(list))
	for _, v := range list {
		items = append(items, toTeacherPageVO(v))
	}
	return &pb.TeacherPageResponse{
		Total: total,
		Pages: page.CalcPages(total, limit),
		List:  items,
	}, nil
}
