//go:build integration

// 集成测试：依赖本地 MySQL 与 Redis。运行：go test -tags integration ./internal/model/
package model

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func newTestLessonModel(t *testing.T) LearningLessonModel {
	t.Helper()
	dsn := os.Getenv("TJXT_TEST_DSN_LEARNING")
	if dsn == "" {
		dsn = "root:0000@tcp(127.0.0.1:3306)/tj_learning?charset=utf8mb4&parseTime=true&loc=Local"
	}
	conn := sqlx.NewMysql(dsn)
	var one int
	if err := conn.QueryRowCtx(context.Background(), &one, "select 1"); err != nil {
		t.Skipf("mysql unavailable, skip integration test: %v", err)
	}
	cc := cache.CacheConf{{RedisConf: redis.RedisConf{Host: envOrDefaultL("TJXT_TEST_REDIS", "127.0.0.1:6379"), Type: "node"}, Weight: 100}}

	return NewLearningLessonModel(conn, cc)
}

func envOrDefaultL(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// TestGrantCoursesIdempotent 锁定开课幂等语义：
// order.pay 事件重复投递时，用户课程记录不产生重复行。
func TestGrantCoursesIdempotent(t *testing.T) {
	m := newTestLessonModel(t)
	ctx := context.Background()
	const userID = 990001
	const courseID = 991001

	t.Cleanup(func() {
		_, _ = m.(*customLearningLessonModel).ExecNoCacheCtx(context.Background(),
			"DELETE FROM learning_lesson WHERE user_id = ? AND course_id = ?", userID, courseID)
	})

	if err := m.GrantCourses(ctx, userID, []int64{courseID}); err != nil {
		t.Fatalf("first GrantCourses: %v", err)
	}
	if err := m.GrantCourses(ctx, userID, []int64{courseID}); err != nil {
		t.Fatalf("second GrantCourses (idempotent replay): %v", err)
	}

	lesson, err := m.FindByUserCourse(ctx, userID, courseID)
	if err != nil {
		t.Fatalf("lesson should exist after grant: %v", err)
	}
	if lesson.Status != LessonStatusNotStart {
		t.Fatalf("status = %d, want 0 (not start)", lesson.Status)
	}
}

// TestRevokeCoursesAfterGrant 退课撤销：状态置失效且可重复执行。
func TestRevokeCoursesAfterGrant(t *testing.T) {
	m := newTestLessonModel(t)
	ctx := context.Background()
	const userID = 990002
	const courseID = 991002

	t.Cleanup(func() {
		_, _ = m.(*customLearningLessonModel).ExecNoCacheCtx(context.Background(),
			"DELETE FROM learning_lesson WHERE user_id = ? AND course_id = ?", userID, courseID)
	})

	if err := m.GrantCourses(ctx, userID, []int64{courseID}); err != nil {
		t.Fatal(err)
	}
	if err := m.RevokeCourses(ctx, userID, []int64{courseID}); err != nil {
		t.Fatalf("RevokeCourses: %v", err)
	}
	if err := m.RevokeCourses(ctx, userID, []int64{courseID}); err != nil {
		t.Fatalf("RevokeCourses replay (must be idempotent): %v", err)
	}

	lesson, err := m.FindByUserCourse(ctx, userID, courseID)
	if err != nil {
		t.Fatal(err)
	}
	if lesson.Status != LessonStatusExpired {
		t.Fatalf("status = %d, want 3 (expired)", lesson.Status)
	}
}

// TestCommitRecordSectionType 锁定 section_type 语义：
// 视频提交更新最近学习小节；考试提交只刷新时间、不动视频进度。
func TestCommitRecordSectionType(t *testing.T) {
	m := newTestLessonModel(t)
	ctx := context.Background()
	userID := time.Now().UnixMilli()
	courseID := userID + 1
	lessonID := userID + 2

	t.Cleanup(func() {
		_, _ = m.(*customLearningLessonModel).ExecNoCacheCtx(context.Background(),
			"DELETE FROM learning_lesson WHERE user_id = ? AND course_id = ?", userID, courseID)
	})

	// 直接插入一条 lesson（模拟已开课），latest_section_id 初始为 777
	if _, err := m.(*customLearningLessonModel).ExecNoCacheCtx(ctx,
		"INSERT INTO learning_lesson (id, user_id, course_id, status, latest_section_id, latest_learn_time) VALUES (?, ?, ?, 1, 777, NOW())",
		lessonID, userID, courseID); err != nil {
		t.Fatal(err)
	}

	// 考试提交（sectionType=2）：latest_section_id 保持 777 不变
	if err := m.UpdateLatestLearnTime(ctx, lessonID, 100, 60); err != nil {
		t.Fatal(err)
	}
	var got LearningLesson
	err := m.(*customLearningLessonModel).QueryRowNoCacheCtx(ctx, &got,
		"SELECT "+learningLessonRows+" FROM "+m.(*customLearningLessonModel).table+" WHERE id = ?", lessonID)
	if err != nil {
		t.Fatal(err)
	}
	if got.LatestSectionId.Int64 != 777 {
		t.Fatalf("exam commit must not touch latest_section_id, got %d", got.LatestSectionId.Int64)
	}

	// 视频提交（sectionType=1）：latest_section_id 更新为 888
	if err := m.UpdateLatestLearn(ctx, lessonID, 888, 100, 60); err != nil {
		t.Fatal(err)
	}
	err = m.(*customLearningLessonModel).QueryRowNoCacheCtx(ctx, &got,
		"SELECT "+learningLessonRows+" FROM "+m.(*customLearningLessonModel).table+" WHERE id = ?", lessonID)
	if err != nil {
		t.Fatal(err)
	}
	if got.LatestSectionId.Int64 != 888 {
		t.Fatalf("video commit must update latest_section_id, got %d", got.LatestSectionId.Int64)
	}
}
