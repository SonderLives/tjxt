package logic

import (
	"context"
	"database/sql"
	"strings"

	"tjxt/apps/exam/rpc/internal/model"
	"tjxt/apps/exam/rpc/internal/svc"
	"tjxt/apps/exam/rpc/pb"
	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveQuestionLogic {
	return &SaveQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 题目管理
func (l *SaveQuestionLogic) SaveQuestion(in *pb.QuestionSaveReq) (*pb.IdReply, error) {
	// 校验参数
	if strings.TrimSpace(in.Name) == "" {
		return nil, xerr.BadRequestf("题干不能为空")
	}
	if in.Type < QuestionTypeSingle || in.Type > QuestionTypeSubjective {
		return nil, xerr.BadRequestf("题目类型非法")
	}
	if in.Difficulty < QuestionDifficultyEasy || in.Difficulty > QuestionDifficultyHard {
		return nil, xerr.BadRequestf("难易度非法")
	}
	if in.Score <= 0 {
		return nil, xerr.BadRequestf("分值必须大于0")
	}

	// 获取登录用户，作为创建人/更新人（上游 API 通过 metadata 透传）
	userId := userIdFromCtx(l.ctx)
	if userId <= 0 {
		userId = 1
	}
	options := sql.NullString{String: in.Options, Valid: in.Options != ""}

	// 新建题目
	if in.Id == 0 {
		id := idgen.NextID()
		q := &model.Question{
			Id:           id,
			Name:         in.Name,
			Type:         int64(in.Type),
			CateId1:      in.CateId1,
			CateId2:      in.CateId2,
			CateId3:      in.CateId3,
			Difficulty:   int64(in.Difficulty),
			AnswerTimes:  0,
			CorrectTimes: 0,
			Score:        int64(in.Score),
			Creater:      userId,
			Updater:      userId,
		}
		if _, err := l.svcCtx.QuestionModel.Insert(l.ctx, q); err != nil {
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "保存题目失败")
		}
		d := &model.QuestionDetail{
			Id:       id,
			Options:  options,
			Answer:   in.Answer,
			Analysis: in.Analysis,
		}
		if _, err := l.svcCtx.QuestionDetailModel.Insert(l.ctx, d); err != nil {
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "保存题目详情失败")
		}
		return &pb.IdReply{Id: id}, nil
	}

	// 更新题目
	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("题目不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询题目失败")
	}
	q.Name = in.Name
	q.Type = int64(in.Type)
	q.CateId1 = in.CateId1
	q.CateId2 = in.CateId2
	q.CateId3 = in.CateId3
	q.Difficulty = int64(in.Difficulty)
	q.Score = int64(in.Score)
	q.Updater = userId
	if err := l.svcCtx.QuestionModel.Update(l.ctx, q); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新题目失败")
	}

	// upsert 题目详情
	d, err := l.svcCtx.QuestionDetailModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if !isNotFound(err) {
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询题目详情失败")
		}
		if _, err := l.svcCtx.QuestionDetailModel.Insert(l.ctx, &model.QuestionDetail{
			Id:       in.Id,
			Options:  options,
			Answer:   in.Answer,
			Analysis: in.Analysis,
		}); err != nil {
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "保存题目详情失败")
		}
	} else {
		d.Options = options
		d.Answer = in.Answer
		d.Analysis = in.Analysis
		if err := l.svcCtx.QuestionDetailModel.Update(l.ctx, d); err != nil {
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新题目详情失败")
		}
	}
	return &pb.IdReply{Id: in.Id}, nil
}
