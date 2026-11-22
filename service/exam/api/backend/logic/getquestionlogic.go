package examlogic

import (
	"context"
	"errors"

	"tjxt/service/exam/api/backend/model"
	"tjxt/service/exam/api/backend/pb"
	"tjxt/service/exam/api/backend/svc"
	"tjxt/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionLogic {
	return &GetQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQuestionLogic) GetQuestion(in *pb.IdReq) (*pb.QuestionVO, error) {
	if in.Id <= 0 {
		return nil, errorx.BadRequestf("题目id非法")
	}

	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("题目不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询题目失败")
	}
	d, err := loadDetail(l.ctx, l.svcCtx.QuestionDetailModel, in.Id)
	if err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询题目详情失败")
	}
	return toQuestionVO(q, d), nil
}
