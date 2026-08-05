package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoursePortalQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCoursePortalQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoursePortalQueryLogic {
	return &CoursePortalQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CoursePortalQuery 用户端课程分页查询，只展示已上架课程。
func (l *CoursePortalQueryLogic) CoursePortalQuery(in *pb.CoursePortalQueryRequest) (*pb.CoursePageQueryReply, error) {
	offset, limit := page.Normalize(in.PageNo, in.PageSize)
	filter := model.CoursePageFilter{
		Keyword:      in.Keyword,
		Status:       CourseStatusUpShelf, // 用户端固定只查已上架课程
		Free:         int64(in.Free),
		CourseType:   int64(in.Type),
		FirstCateId:  in.CategoryIdLv1,
		SecondCateId: in.CategoryIdLv2,
		ThirdCateId:  in.CategoryIdLv3,
		BeginTime:    in.BeginTime,
		EndTime:      in.EndTime,
	}

	list, total, err := l.svcCtx.CourseModel.PageQuery(l.ctx, filter, offset, limit)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "分页查询课程失败")
	}

	return &pb.CoursePageQueryReply{
		Total: total,
		Pages: page.CalcPages(total, limit),
		List:  l.buildItems(list),
	}, nil
}

// buildItems 课程实体转分页 VO，并一次性回填一/二/三级分类名称。
func (l *CoursePortalQueryLogic) buildItems(list []*model.Course) []*pb.CoursePageItem {
	items := make([]*pb.CoursePageItem, 0, len(list))
	if len(list) == 0 {
		return items
	}
	cateMap := map[int64]*model.Category{}
	if all, err := l.svcCtx.CategoryModel.ListAll(l.ctx); err == nil {
		cateMap = categoryNameMap(all)
	}
	for _, c := range list {
		item := &pb.CoursePageItem{
			Id:            c.Id,
			Name:          c.Name,
			CoverUrl:      c.CoverUrl,
			Price:         c.Price,
			Free:          int32(c.Free),
			Status:        int32(c.Status),
			StatusDesc:    courseStatusDesc(c.Status),
			Score:         c.Score,
			Sold:          0, // 销量来自 trade 服务，course 库无对应列
			SectionNum:    formatNullInt64(c.SectionNum),
			MediaDuration: formatNullInt64(c.MediaDuration),
			PublishTime:   formatNullTime(c.PublishTime),
			CreateTime:    formatTime(c.CreateTime),
			UpdateTime:    formatTime(c.UpdateTime),
		}
		if cate := cateMap[c.FirstCateId]; cate != nil {
			item.FirstCateName = cate.Name
		}
		if cate := cateMap[c.SecondCateId]; cate != nil {
			item.SecondCateName = cate.Name
		}
		if cate := cateMap[c.ThirdCateId]; cate != nil {
			item.ThirdCateName = cate.Name
		}
		items = append(items, item)
	}
	return items
}
