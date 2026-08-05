package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSimpleInfoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSimpleInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSimpleInfoListLogic {
	return &CourseSimpleInfoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseSimpleInfoList 按课程 id 或三级分类 id 批量查询课程简要信息，供其他服务调用。
func (l *CourseSimpleInfoListLogic) CourseSimpleInfoList(in *pb.CourseSimpleInfoQueryRequest) (*pb.CourseSimpleInfoListReply, error) {
	var (
		list []*model.Course
		err  error
	)
	switch {
	case len(in.Ids) > 0:
		list, err = l.svcCtx.CourseModel.ListByIds(l.ctx, in.Ids)
	case len(in.ThirdCateIds) > 0:
		list, err = l.svcCtx.CourseModel.ListByThirdCateIds(l.ctx, in.ThirdCateIds)
	default:
		return &pb.CourseSimpleInfoListReply{Items: []*pb.CourseSimpleInfoItem{}}, nil
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "批量查询课程失败")
	}

	items := make([]*pb.CourseSimpleInfoItem, 0, len(list))
	for _, c := range list {
		items = append(items, &pb.CourseSimpleInfoItem{
			Id:              c.Id,
			Name:            c.Name,
			CoverUrl:        c.CoverUrl,
			Price:           c.Price,
			Free:            int32(c.Free),
			SectionNum:      formatNullInt64(c.SectionNum),
			Status:          int32(c.Status),
			ValidDuration:   c.ValidDuration,
			PurchaseEndTime: formatTime(c.PurchaseEndTime),
			FirstCateId:     c.FirstCateId,
			SecondCateId:    c.SecondCateId,
			ThirdCateId:     c.ThirdCateId,
		})
	}
	return &pb.CourseSimpleInfoListReply{Items: items}, nil
}
