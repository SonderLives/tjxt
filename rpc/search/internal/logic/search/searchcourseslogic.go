package searchlogic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"tjxt/apps/search/rpc/internal/svc"
	"tjxt/apps/search/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCoursesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCoursesLogic {
	return &SearchCoursesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SearchCourses 课程全文检索：match 查 name 字段 + enable=1 过滤，支持分页与排序。
// Sort：0 默认相关度，1 按销量，2 按评分。
func (l *SearchCoursesLogic) SearchCourses(in *pb.CourseSearchRequest) (*pb.CourseSearchPageReply, error) {
	if strings.TrimSpace(in.Keyword) == "" {
		return nil, xerr.BadRequestf("搜索关键词不能为空")
	}
	if l.svcCtx.ES == nil {
		return nil, xerr.ServiceUnavailable("搜索服务暂不可用")
	}

	pageNo := in.PageNo
	if pageNo < 1 {
		pageNo = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	sort := ""
	switch in.Sort {
	case 1:
		sort = `, "sort": [{"sold": {"order": "desc"}}]`
	case 2:
		sort = `, "sort": [{"score": {"order": "desc"}}]`
	}

	body := fmt.Sprintf(`{
  "query": {
    "bool": {
      "must": [
        {"match": {"name": %s}}
      ],
      "filter": [
        {"term": {"enable": 1}}
      ]
    }
  },
  "from": %d,
  "size": %d%s
}`, strconv.Quote(in.Keyword), (pageNo-1)*pageSize, pageSize, sort)

	return courseSearch(l.ctx, l.svcCtx, []byte(body))
}
