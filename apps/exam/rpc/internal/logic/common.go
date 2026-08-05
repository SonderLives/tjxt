package logic

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"tjxt/apps/exam/rpc/internal/model"
	"tjxt/apps/exam/rpc/pb"

	"google.golang.org/grpc/metadata"
)

// 题目类型
const (
	QuestionTypeSingle     = 1 // 单选题
	QuestionTypeMultiple   = 2 // 多选题
	QuestionTypeUncertain  = 3 // 不定向选择题
	QuestionTypeJudge      = 4 // 判断题
	QuestionTypeSubjective = 5 // 主观题
)

// 难易度
const (
	QuestionDifficultyEasy   = 1 // 简单
	QuestionDifficultyMedium = 2 // 中等
	QuestionDifficultyHard   = 3 // 困难
)

// metadataKeyUserID 上游 API 通过 gRPC metadata 透传的当前登录用户 id
const metadataKeyUserID = "user_id"

// isNotFound 判断错误是否是数据不存在
func isNotFound(err error) bool {
	return err == sql.ErrNoRows || err == model.ErrNotFound
}

// userIdFromCtx 从 gRPC 入站 metadata 中读取上游透传的登录用户 id
func userIdFromCtx(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}
	vals := md.Get(metadataKeyUserID)
	if len(vals) == 0 {
		return 0
	}
	id, err := strconv.ParseInt(vals[0], 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// formatTime 时间统一格式化为 yyyy-MM-dd HH:mm:ss
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// toQuestionVO 组装题目完整 VO（含详情）
func toQuestionVO(q *model.Question, d *model.QuestionDetail) *pb.QuestionVO {
	vo := &pb.QuestionVO{
		Id:           q.Id,
		Name:         q.Name,
		Type:         int32(q.Type),
		CateId1:      q.CateId1,
		CateId2:      q.CateId2,
		CateId3:      q.CateId3,
		Difficulty:   int32(q.Difficulty),
		AnswerTimes:  int32(q.AnswerTimes),
		CorrectTimes: int32(q.CorrectTimes),
		Score:        int32(q.Score),
		CreateTime:   formatTime(q.CreateTime),
	}
	if d != nil {
		vo.Options = d.Options.String
		vo.Answer = d.Answer
		vo.Analysis = d.Analysis
	}
	return vo
}

// loadDetail 加载题目详情，不存在时返回 nil
func loadDetail(ctx context.Context, detailModel model.QuestionDetailModel, id int64) (*model.QuestionDetail, error) {
	d, err := detailModel.FindOne(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}

// loadQuestionsByIds 按 id 集合批量加载题目，返回 id -> 题目 的映射
func loadQuestionsByIds(ctx context.Context, questionModel model.QuestionModel, ids []int64) (map[int64]*model.Question, error) {
	list, err := questionModel.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*model.Question, len(list))
	for _, q := range list {
		m[q.Id] = q
	}
	return m, nil
}
