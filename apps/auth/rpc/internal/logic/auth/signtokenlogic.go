package authlogic

import (
	"context"
	"time"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	authutil "tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSignTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignTokenLogic {
	return &SignTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SignTokenLogic) SignToken(in *pb.SignTokenReq) (*pb.SignTokenReply, error) {
	expire := in.ExpireSec
	if expire <= 0 {
		expire = l.svcCtx.Config.Jwt.AccessExpire
	}
	if expire <= 0 {
		expire = 7200
	}
	access, err := authutil.Sign(in.UserId, in.RoleCode, l.svcCtx.Config.Jwt.AccessSecret, expire)
	if err != nil {
		return nil, err
	}
	refresh, err := authutil.Sign(in.UserId, in.RoleCode, l.svcCtx.Config.Jwt.RefreshSecret, l.svcCtx.Config.Jwt.RefreshExpire)
	if err != nil {
		return nil, err
	}
	return &pb.SignTokenReply{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    time.Now().Unix() + expire,
	}, nil
}