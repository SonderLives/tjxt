package logic

import (
	"context"

	"tjxt/apps/exam/rpc/internal/svc"
	"tjxt/apps/exam/rpc/pb"
	"tjxt/pkg/xerr"

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
		return nil, xerr.BadRequestf("题目id非法")
	}

	q, err := l.svcCtx.QuestionModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("题目不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询题目失败")
	}
	d, err := loadDetail(l.ctx, l.svcCtx.QuestionDetailModel, in.Id)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询题目详情失败")
	}
	return toQuestionVO(q, d), nil
}
