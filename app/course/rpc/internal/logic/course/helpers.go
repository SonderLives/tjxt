package courselogic

import (
	"context"

	"tjxt/app/course/rpc/internal/svc"
	userclient "tjxt/app/user/rpc/client/user"

	"github.com/zeromicro/go-zero/core/logx"
)

// teacherInfo 教师展示信息（姓名/头像/职位/简介来自 user 服务）。
type teacherInfo struct {
	Name      string
	Icon      string
	Photo     string
	Job       string
	Introduce string
}

// fetchTeacherMap 批量查询教师（user 服务），返回 user_id -> 展示信息。
// user 服务不可用时返回空 map，调用方回落到空字符串展示，不阻塞课程主流程。
func fetchTeacherMap(ctx context.Context, svcCtx *svc.ServiceContext, teacherIds []int64) map[int64]*teacherInfo {
	res := make(map[int64]*teacherInfo, len(teacherIds))
	if len(teacherIds) == 0 {
		return res
	}
	reply, err := svcCtx.UserRpc.GetUsersByIds(ctx, &userclient.UserIdsRequest{UserIds: teacherIds})
	if err != nil {
		logx.WithContext(ctx).Errorf("query teachers from user service failed, ids=%v: %v", teacherIds, err)
		return res
	}
	for _, u := range reply.List {
		res[u.Id] = &teacherInfo{Name: u.Name, Icon: u.Icon, Photo: u.Icon, Job: u.Job, Introduce: u.Intro}
	}
	return res
}
