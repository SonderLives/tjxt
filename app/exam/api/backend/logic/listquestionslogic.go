package examlogic

import (
	"context"

	"tjxt/app/exam/api/backend/pb"
	"tjxt/app/exam/api/backend/svc"
	"tjxt/pkg/errorx"
	"tjxt/pkg/page"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListQuestionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListQuestionsLogic {
	return &ListQuestionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListQuestionsLogic) ListQuestions(in *pb.QuestionListReq) (*pb.QuestionListReply, error) {
	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))
	list, total, err := l.svcCtx.QuestionModel.FindPage(l.ctx,
		in.Name, int64(in.Type), in.CateId1, in.CateId2, int64(in.Difficulty), offset, limit)
	if err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "分页查询题目失败")
	}

	resp := &pb.QuestionListReply{
		Total: total,
		List:  make([]*pb.QuestionVO, 0, len(list)),
	}
	for _, q := range list {
		d, err := loadDetail(l.ctx, l.svcCtx.QuestionDetailModel, q.Id)
		if err != nil {
			return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询题目详情失败")
		}
		resp.List = append(resp.List, toQuestionVO(q, d))
	}
	return resp, nil
}
