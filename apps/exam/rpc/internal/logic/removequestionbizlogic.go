package logic

import (
	"context"
	"database/sql"

	"tjxt/apps/exam/rpc/internal/svc"
	"tjxt/apps/exam/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveQuestionBizLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveQuestionBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveQuestionBizLogic {
	return &RemoveQuestionBizLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveQuestionBizLogic) RemoveQuestionBiz(in *pb.QuestionBizReq) (*pb.Empty, error) {
	if in.BizId <= 0 {
		return nil, xerr.BadRequestf("业务id非法")
	}
	if in.QuestionId <= 0 {
		return nil, xerr.BadRequestf("题目id非法")
	}

	bizId := sql.NullInt64{Int64: in.BizId, Valid: true}
	questionId := sql.NullInt64{Int64: in.QuestionId, Valid: true}

	// 幂等：未绑定也返回成功
	existing, err := l.svcCtx.QuestionBizModel.FindOneByBizIdQuestionId(l.ctx, bizId, questionId)
	if err != nil {
		if isNotFound(err) {
			return &pb.Empty{}, nil
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询题目业务关联失败")
	}
	if err := l.svcCtx.QuestionBizModel.Delete(l.ctx, existing.Id); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "解除题目绑定失败")
	}
	return &pb.Empty{}, nil
}
