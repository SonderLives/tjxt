package logic

import (
	"context"

	courseclient "tjxt/apps/course/rpc/course"
	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

// enrichLessons 用 CourseRpc 回填课程维度的字段：名称/封面/价格/小节数。
// learning 自身不持有课程信息，这些字段来自 course 服务。
func enrichLessons(ctx context.Context, svcCtx *svc.ServiceContext, vos []*pb.LearningLessonVO) {
	if len(vos) == 0 {
		return
	}
	ids := make([]int64, 0, len(vos))
	idSet := make(map[int64]struct{}, len(vos))
	for _, v := range vos {
		if _, ok := idSet[v.CourseId]; ok {
			continue
		}
		idSet[v.CourseId] = struct{}{}
		ids = append(ids, v.CourseId)
	}
	reply, err := svcCtx.CourseRpc.CourseSimpleInfoList(ctx, &courseclient.CourseSimpleInfoQueryRequest{Ids: ids})
	if err != nil {
		logx.WithContext(ctx).Errorf("enrichLessons CourseSimpleInfoList failed: %v", err)
		return
	}
	infoMap := make(map[int64]*courseclient.CourseSimpleInfoItem, len(reply.Items))
	for _, it := range reply.Items {
		infoMap[it.Id] = it
	}
	for _, v := range vos {
		if it, ok := infoMap[v.CourseId]; ok {
			v.CourseName = it.Name
			v.CourseCoverUrl = it.CoverUrl
			v.CourseAmount = it.Price
			v.Sections = int32(it.SectionNum)
		}
	}
}

// toLessonVOTypes 将 pb 视图对象转为 API 对外 types。
func toLessonVOTypes(v *pb.LearningLessonVO) types.LearningLessonVO {
	if v == nil {
		return types.LearningLessonVO{}
	}
	return types.LearningLessonVO{
		Id:                 v.Id,
		CourseId:           v.CourseId,
		CourseName:         v.CourseName,
		CourseCoverUrl:     v.CourseCoverUrl,
		CourseAmount:       v.CourseAmount,
		Sections:           int64(v.Sections),
		LearnedSections:    int64(v.LearnedSections),
		Status:             v.Status,
		PlanStatus:         v.PlanStatus,
		WeekFreq:           int64(v.WeekFreq),
		LatestSectionId:    v.LatestSectionId,
		LatestSectionName:  v.LatestSectionName,
		LatestSectionIndex: int64(v.LatestSectionIndex),
		CreateTime:         v.CreateTime,
		ExpireTime:         v.ExpireTime,
		LatestLearnTime:    v.LatestLearnTime,
	}
}
