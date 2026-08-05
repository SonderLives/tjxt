package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseMediaUseInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseMediaUseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseMediaUseInfoLogic {
	return &CourseMediaUseInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseMediaUseInfo 统计媒资在课程目录中的引用次数，供媒资服务判断能否删除。
func (l *CourseMediaUseInfoLogic) CourseMediaUseInfo(in *pb.MediaIdsRequest) (*pb.MediaQuoteList, error) {
	if len(in.MediaIds) == 0 {
		return &pb.MediaQuoteList{Items: []*pb.MediaQuote{}}, nil
	}
	counts, err := l.svcCtx.CourseCatalogueModel.CountByMediaIds(l.ctx, in.MediaIds)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "统计媒资引用次数失败")
	}
	quoteMap := make(map[int64]int64, len(counts))
	for _, c := range counts {
		quoteMap[c.MediaId] = c.QuoteNum
	}
	// 未被引用的媒资也返回，引用数为 0
	items := make([]*pb.MediaQuote, 0, len(in.MediaIds))
	for _, mediaId := range in.MediaIds {
		items = append(items, &pb.MediaQuote{
			MediaId:  mediaId,
			QuoteNum: quoteMap[mediaId],
		})
	}
	return &pb.MediaQuoteList{Items: items}, nil
}
