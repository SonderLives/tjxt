package logic

import (
	"learning/internal/model"
	"learning/rpc/pb/pb"
)

func lessonReply(lesson *model.LearningLesson) *pb.LessonReply {
	reply := &pb.LessonReply{Id: lesson.Id, CourseId: lesson.CourseId, Status: int32(lesson.Status), PlanStatus: int32(lesson.PlanStatus), LearnedSections: int32(lesson.LearnedSections), CreateTime: lesson.CreateTime.Unix()}
	if lesson.WeekFreq.Valid {
		reply.WeekFreq = int32(lesson.WeekFreq.Int64)
	}
	if lesson.LatestSectionId.Valid {
		reply.LatestSectionId = lesson.LatestSectionId.Int64
	}
	if lesson.ExpireTime.Valid {
		reply.ExpireTime = lesson.ExpireTime.Time.Unix()
	}
	return reply
}
