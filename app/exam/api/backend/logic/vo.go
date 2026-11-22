package examlogic

import (
	"context"
	"errors"

	"tjxt/app/exam/api/backend/model"
	"tjxt/app/exam/api/backend/pb"
	"tjxt/pkg/timex"
)

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
		CreateTime:   timex.Format(q.CreateTime),
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
		if errors.Is(err, model.ErrNotFound) {
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
