package logic

import (
	"context"
	"database/sql"

	"tjxt/apps/exam/rpc/internal/model"
	"tjxt/apps/exam/rpc/internal/svc"
	"tjxt/apps/exam/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddQuestionBizLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddQuestionBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddQuestionBizLogic {
	return &AddQuestionBizLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 题目业务关联
func (l *AddQuestionBizLogic) AddQuestionBiz(in *pb.QuestionBizReq) (*pb.IdReply, error) {
	if in.BizId <= 0 {
		return nil, xerr.BadRequestf("业务id非法")
	}
	if in.QuestionId <= 0 {
		return nil, xerr.BadRequestf("题目id非法")
	}

	// 题目必须存在
	if _, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.QuestionId); err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("题目不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询题目失败")
	}

	bizId := sql.NullInt64{Int64: in.BizId, Valid: true}
	questionId := sql.NullInt64{Int64: in.QuestionId, Valid: true}

	// 幂等：已绑定则直接返回
	existing, err := l.svcCtx.QuestionBizModel.FindOneByBizIdQuestionId(l.ctx, bizId, questionId)
	if err == nil {
		return &pb.IdReply{Id: existing.Id}, nil
	}
	if !isNotFound(err) {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询题目业务关联失败")
	}

	ret, err := l.svcCtx.QuestionBizModel.Insert(l.ctx, &model.QuestionBiz{
		BizId:      bizId,
		QuestionId: questionId,
	})
	if err != nil {
		// 并发下唯一索引冲突，同样按幂等处理
		existing, err2 := l.svcCtx.QuestionBizModel.FindOneByBizIdQuestionId(l.ctx, bizId, questionId)
		if err2 == nil {
			return &pb.IdReply{Id: existing.Id}, nil
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "绑定题目失败")
	}
	id, err := ret.LastInsertId()
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "获取绑定记录id失败")
	}
	return &pb.IdReply{Id: id}, nil
}
