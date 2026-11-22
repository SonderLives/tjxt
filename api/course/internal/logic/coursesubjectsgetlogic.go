// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSubjectsGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSubjectsGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSubjectsGetLogic {
	return &CourseSubjectsGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseSubjectsGet 查询课程各目录下的题目（透传 RPC）。
func (l *CourseSubjectsGetLogic) CourseSubjectsGet(req *types.IdPathReq) (resp []types.CataSubjectVO, err error) {
	list, err := l.svcCtx.CourseRpc.CourseSubjectsGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程题目失败")
	}
	resp = make([]types.CataSubjectVO, 0, len(list.Items))
	for _, c := range list.Items {
		subjects := make([]types.SubjectInfoVO, 0, len(c.Subjects))
		for _, s := range c.Subjects {
			subjects = append(subjects, types.SubjectInfoVO{
				Id:          s.Id,
				Name:        s.Name,
				SubjectType: int64(s.SubjectType),
				Difficulty:  int64(s.Difficulty),
				Options:     s.Options,
				Answer:      s.Answer,
				Analysis:    s.Analysis,
				Score:       s.Score,
			})
		}
		resp = append(resp, types.CataSubjectVO{
			CataId:   c.CataId,
			Subjects: subjects,
		})
	}
	return resp, nil
}
