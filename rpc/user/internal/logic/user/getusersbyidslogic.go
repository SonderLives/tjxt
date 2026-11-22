package userlogic

import (
	"context"

	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsersByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersByIdsLogic {
	return &GetUsersByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUsersByIds 批量获取用户（含资料），空入参直接返回空列表。
func (l *GetUsersByIdsLogic) GetUsersByIds(in *pb.UserIdsRequest) (*pb.UserListResponse, error) {
	if len(in.UserIds) == 0 {
		return &pb.UserListResponse{}, nil
	}
	list, err := l.svcCtx.UserModel.FindByIdsWithDetail(l.ctx, in.UserIds)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询用户失败")
	}
	dtos := make([]*pb.UserDTO, 0, len(list))
	for _, v := range list {
		dtos = append(dtos, toUserDTOFromView(v))
	}
	return &pb.UserListResponse{List: dtos}, nil
}
