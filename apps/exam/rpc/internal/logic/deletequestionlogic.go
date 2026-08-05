package logic

import (
	"context"

	"tjxt/apps/exam/rpc/internal/svc"
	"tjxt/apps/exam/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionLogic {
	return &DeleteQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteQuestionLogic) DeleteQuestion(in *pb.IdReq) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("题目id非法")
	}

	// 删除题目（缓存由 CachedConn 自动清理，幂等：不存在同样返回成功）
	if err := l.svcCtx.QuestionModel.Delete(l.ctx, in.Id); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "删除题目失败")
	}
	// 删除题目详情
	if err := l.svcCtx.QuestionDetailModel.Delete(l.ctx, in.Id); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "删除题目详情失败")
	}
	// 删除该题目的所有业务关联
	if err := l.svcCtx.QuestionBizModel.DeleteByQuestionId(l.ctx, in.Id); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "删除题目业务关联失败")
	}
	return &pb.Empty{}, nil
}
