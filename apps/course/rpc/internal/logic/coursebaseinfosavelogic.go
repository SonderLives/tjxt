package logic

import (
	"context"
	"database/sql"
	"time"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseBaseInfoSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseBaseInfoSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseBaseInfoSaveLogic {
	return &CourseBaseInfoSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseBaseInfoSave 新增或更新课程草稿基础信息，同时保存课程内容（介绍/详情/适用人群）。
func (l *CourseBaseInfoSaveLogic) CourseBaseInfoSave(in *pb.CourseBaseInfoSaveRequest) (*pb.IdResponse, error) {
	if in.Name == "" {
		return nil, xerr.BadRequestf("课程名称不能为空")
	}

	startTime := l.parseTime(in.PurchaseStartTime)
	endTime := l.parseTime(in.PurchaseEndTime)

	var id int64
	if in.Id == 0 {
		id = nextID()
		draft := &model.CourseDraft{
			Id:                id,
			Name:              in.Name,
			CoverUrl:          in.CoverUrl,
			Price:             in.Price,
			Free:              int64(in.Free),
			ThirdCateId:       in.ThirdCateId,
			ValidDuration:     in.ValidDuration,
			PurchaseStartTime: sql.NullTime{Time: startTime, Valid: !startTime.IsZero()},
			PurchaseEndTime:   endTime,
			Status:            CourseStatusPending,
			Step:              1,
			CanUpdate:         1,
			DepId:             0,
			Creater:           0,
			Updater:           0,
		}
		if _, err := l.svcCtx.CourseDraftModel.Insert(l.ctx, draft); err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "新增课程草稿失败")
		}
	} else {
		id = in.Id
		draft, err := l.svcCtx.CourseDraftModel.FindOne(l.ctx, id)
		if err != nil {
			if isNotFound(err) {
				return nil, xerr.NotFound("课程不存在")
			}
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
		}
		draft.Name = in.Name
		draft.CoverUrl = in.CoverUrl
		draft.Price = in.Price
		draft.Free = int64(in.Free)
		draft.ThirdCateId = in.ThirdCateId
		draft.ValidDuration = in.ValidDuration
		draft.PurchaseStartTime = sql.NullTime{Time: startTime, Valid: !startTime.IsZero()}
		draft.PurchaseEndTime = endTime
		draft.Updater = 0
		if draft.Step < 1 {
			draft.Step = 1
		}
		if err := l.svcCtx.CourseDraftModel.Update(l.ctx, draft); err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "更新课程草稿失败")
		}
	}

	now := time.Now()
	content := &model.CourseContentDraft{
		Id:              id,
		CourseIntroduce: in.Introduce,
		CourseDetail:    in.Detail,
		UsePeople:       in.UsePeople,
		DepId:           0,
		CreateTime:      now,
		UpdateTime:      now,
		Creater:         0,
		Updater:         0,
	}
	if err := l.svcCtx.CourseContentDraftModel.Upsert(l.ctx, content); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "保存课程内容失败")
	}
	return &pb.IdResponse{Id: id}, nil
}

// parseTime 解析前端传入的时间字符串，非法或空值返回零值时间。
func (l *CourseBaseInfoSaveLogic) parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t
		}
	}
	return time.Time{}
}
