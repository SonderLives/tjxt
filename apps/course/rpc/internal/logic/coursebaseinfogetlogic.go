package logic

import (
	"context"
	"fmt"
	"strings"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseBaseInfoGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseBaseInfoGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseBaseInfoGetLogic {
	return &CourseBaseInfoGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 课程基础信息 =====
// CourseBaseInfoGet 查询课程草稿的基础信息（编辑态），附带课程内容、目录数量与分类名称。
func (l *CourseBaseInfoGetLogic) CourseBaseInfoGet(in *pb.IdRequest) (*pb.CourseBaseInfoView, error) {
	if in.Id == 0 {
		return nil, xerr.BadRequestf("课程id不能为空")
	}
	d, err := l.svcCtx.CourseDraftModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("课程不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}

	view := &pb.CourseBaseInfoView{
		Id:                d.Id,
		Name:              d.Name,
		Price:             d.Price,
		Free:              int32(d.Free),
		CoverUrl:          d.CoverUrl,
		ValidDuration:     d.ValidDuration,
		PurchaseStartTime: formatNullTime(d.PurchaseStartTime),
		PurchaseEndTime:   formatTime(d.PurchaseEndTime),
		FirstCateId:       d.FirstCateId,
		SecondCateId:      d.SecondCateId,
		ThirdCateId:       d.ThirdCateId,
		Status:            int32(d.Status),
		Step:              int32(d.Step),
		Score:             d.Score,
		CanUpdate:         d.CanUpdate == 1,
		CreateTime:        formatTime(d.CreateTime),
		UpdateTime:        formatTime(d.UpdateTime),
		// 报名数/学习人数/实付金额/退款数、创建人与更新人姓名来自其他服务，course 库无对应数据
		EnrollNum:     0,
		StudyNum:      0,
		RealPayAmount: 0,
		RefundNum:     0,
		CreaterName:   "",
		UpdaterName:   "",
	}
	if d.Score > 0 {
		view.CourseScore = fmt.Sprintf("%.1f", float64(d.Score)/10)
	}

	// 课程内容（介绍/详情/适用人群）
	content, cerr := l.svcCtx.CourseContentDraftModel.FindOne(l.ctx, d.Id)
	if cerr != nil && !isNotFound(cerr) {
		return nil, xerr.Wrap(cerr, xerr.CodeInternal, "查询课程内容失败")
	}
	if content != nil {
		view.Introduce = content.CourseIntroduce
		view.Detail = content.CourseDetail
		view.UsePeople = content.UsePeople
	}

	// 目录总数
	catalogues, kerr := l.svcCtx.CourseCatalogueDraftModel.ListByCourseId(l.ctx, d.Id)
	if kerr != nil {
		return nil, xerr.Wrap(kerr, xerr.CodeInternal, "查询课程目录失败")
	}
	view.CataTotalNum = int64(len(catalogues))

	// 一/二/三级分类名称拼接
	if all, aerr := l.svcCtx.CategoryModel.ListAll(l.ctx); aerr == nil {
		m := categoryNameMap(all)
		names := make([]string, 0, 3)
		for _, cateId := range []int64{d.FirstCateId, d.SecondCateId, d.ThirdCateId} {
			if cateId == 0 {
				continue
			}
			if c := m[cateId]; c != nil {
				names = append(names, c.Name)
			}
		}
		view.CateNames = strings.Join(names, "/")
	}
	return view, nil
}
