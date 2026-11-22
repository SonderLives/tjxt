package examlogic

import (
	"context"

	"tjxt/rpc/exam/internal/svc"
	"tjxt/rpc/exam/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionsByBizLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionsByBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionsByBizLogic {
	return &GetQuestionsByBizLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQuestionsByBizLogic) GetQuestionsByBiz(in *pb.QuestionBizListReq) (*pb.QuestionListReply, error) {
	if in.BizId <= 0 {
		return nil, errorx.BadRequestf("业务id非法")
	}

	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))
	// 先分页查出该业务下的题目关联
	bizList, total, err := l.svcCtx.QuestionBizModel.FindPageByBizId(l.ctx, in.BizId, offset, limit)
	if err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "分页查询题目业务关联失败")
	}

	ids := make([]int64, 0, len(bizList))
	for _, b := range bizList {
		ids = append(ids, b.QuestionId.Int64)
	}
	questionMap, err := loadQuestionsByIds(l.ctx, l.svcCtx.QuestionModel, ids)
	if err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询题目失败")
	}

	resp := &pb.QuestionListReply{
		Total: total,
		List:  make([]*pb.QuestionVO, 0, len(bizList)),
	}
	for _, b := range bizList {
		q, ok := questionMap[b.QuestionId.Int64]
		if !ok {
			// 关联存在但题目已删除，跳过
			continue
		}
		d, err := loadDetail(l.ctx, l.svcCtx.QuestionDetailModel, q.Id)
		if err != nil {
			return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询题目详情失败")
		}
		resp.List = append(resp.List, toQuestionVO(q, d))
	}
	return resp, nil
}
