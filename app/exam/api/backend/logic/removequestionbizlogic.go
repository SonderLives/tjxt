package examlogic

import (
	"context"
	"database/sql"
	"errors"

	"tjxt/app/exam/api/backend/model"
	"tjxt/app/exam/api/backend/pb"
	"tjxt/app/exam/api/backend/svc"
	"tjxt/pkg/errorx"

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
		return nil, errorx.BadRequestf("业务id非法")
	}
	if in.QuestionId <= 0 {
		return nil, errorx.BadRequestf("题目id非法")
	}

	bizId := sql.NullInt64{Int64: in.BizId, Valid: true}
	questionId := sql.NullInt64{Int64: in.QuestionId, Valid: true}

	// 幂等：未绑定也返回成功
	existing, err := l.svcCtx.QuestionBizModel.FindOneByBizIdQuestionId(l.ctx, bizId, questionId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &pb.Empty{}, nil
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询题目业务关联失败")
	}
	if err := l.svcCtx.QuestionBizModel.Delete(l.ctx, existing.Id); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "解除题目绑定失败")
	}
	return &pb.Empty{}, nil
}
