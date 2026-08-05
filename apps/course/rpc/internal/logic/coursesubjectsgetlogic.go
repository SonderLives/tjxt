package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSubjectsGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSubjectsGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSubjectsGetLogic {
	return &CourseSubjectsGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 章节/小节的题目 =====
// 按课程 id 查询草稿态的「小节-题目」绑定关系，按 cata_id 分组返回题目详情。
func (l *CourseSubjectsGetLogic) CourseSubjectsGet(in *pb.IdRequest) (*pb.CataSubjectInfoList, error) {
	binds, err := l.svcCtx.CourseCataSubjectDraftModel.ListByCourseId(l.ctx, in.Id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程题目关系失败")
	}
	if len(binds) == 0 {
		return &pb.CataSubjectInfoList{Items: []*pb.CatalogSubjectInfo{}}, nil
	}

	// 按 cata_id 分组，保持首次出现顺序
	cataIds := make([]int64, 0, len(binds))
	group := make(map[int64][]int64, len(binds))
	subjectIds := make([]int64, 0, len(binds))
	seenSubject := make(map[int64]struct{}, len(binds))
	for _, b := range binds {
		if _, ok := group[b.CataId]; !ok {
			cataIds = append(cataIds, b.CataId)
		}
		group[b.CataId] = append(group[b.CataId], b.SubjectId)
		if _, ok := seenSubject[b.SubjectId]; !ok {
			seenSubject[b.SubjectId] = struct{}{}
			subjectIds = append(subjectIds, b.SubjectId)
		}
	}

	subjects, err := l.svcCtx.SubjectModel.ListByIds(l.ctx, subjectIds)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询题目详情失败")
	}
	subjectMap := make(map[int64]*model.Subject, len(subjects))
	for _, s := range subjects {
		subjectMap[s.Id] = s
	}

	items := make([]*pb.CatalogSubjectInfo, 0, len(cataIds))
	for _, cataId := range cataIds {
		ids := group[cataId]
		infos := make([]*pb.SubjectInfo, 0, len(ids))
		for _, sid := range ids {
			s, ok := subjectMap[sid]
			if !ok {
				continue
			}
			infos = append(infos, &pb.SubjectInfo{
				Id:          s.Id,
				Name:        s.Name,
				SubjectType: int32(s.SubjectType),
				Difficulty:  int32(s.Difficulty),
				Score:       s.Score,
				Options:     subjectOptions(s),
				Answer:      s.Answer,
				Analysis:    s.Analysis,
			})
		}
		items = append(items, &pb.CatalogSubjectInfo{
			CataId:   cataId,
			Subjects: infos,
		})
	}
	return &pb.CataSubjectInfoList{Items: items}, nil
}

// subjectOptions 收集 option1..option10 中的非空选项。
func subjectOptions(s *model.Subject) []string {
	all := []string{
		s.Option1, s.Option2, s.Option3, s.Option4, s.Option5,
		s.Option6, s.Option7, s.Option8, s.Option9, s.Option10,
	}
	options := make([]string, 0, len(all))
	for _, o := range all {
		if o != "" {
			options = append(options, o)
		}
	}
	return options
}
