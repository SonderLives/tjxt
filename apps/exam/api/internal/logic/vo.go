// Package logic 存放 exam-api 各接口的业务逻辑。
// 本文件提供 pb.VO 到 api types.VO 的转换工具。
package logic

import (
	"tjxt/apps/exam/api/internal/types"
	examclient "tjxt/apps/exam/rpc/exam"
)

// ToQuestionVO 将 RPC 返回的题目 VO 转换为 API 响应 VO
func ToQuestionVO(in *examclient.QuestionVO) *types.QuestionVO {
	if in == nil {
		return nil
	}
	return &types.QuestionVO{
		Id:           in.Id,
		Name:         in.Name,
		Type:         in.Type,
		CateId1:      in.CateId1,
		CateId2:      in.CateId2,
		CateId3:      in.CateId3,
		Difficulty:   in.Difficulty,
		AnswerTimes:  in.AnswerTimes,
		CorrectTimes: in.CorrectTimes,
		Score:        in.Score,
		CreateTime:   in.CreateTime,
		Options:      in.Options,
		Answer:       in.Answer,
		Analysis:     in.Analysis,
	}
}
