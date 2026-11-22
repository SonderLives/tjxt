package searchlogic

import (
	"context"
	"fmt"

	"tjxt/apps/search/rpc/internal/svc"
	"tjxt/apps/search/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTopCoursesByCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTopCoursesByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTopCoursesByCategoryLogic {
	return &GetTopCoursesByCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetTopCoursesByCategory 根据二级分类id查询课程TOP10：
// term 过滤 category_id_lv2 + enable=1，按销量、评分倒序。
func (l *GetTopCoursesByCategoryLogic) GetTopCoursesByCategory(in *pb.IdReq) (*pb.CourseListReply, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("二级分类id非法")
	}
	if l.svcCtx.ES == nil {
		return nil, xerr.ServiceUnavailable("搜索服务暂不可用")
	}

	body := fmt.Sprintf(`{
  "query": {
    "bool": {
      "filter": [
        {"term": {"category_id_lv2": %d}},
        {"term": {"enable": 1}}
      ]
    }
  },
  "sort": [
    {"sold": {"order": "desc"}},
    {"score": {"order": "desc"}}
  ],
  "size": 10
}`, in.Id)

	reply, err := courseSearch(l.ctx, l.svcCtx, []byte(body))
	if err != nil {
		return nil, err
	}
	return &pb.CourseListReply{Items: reply.Items}, nil
}
