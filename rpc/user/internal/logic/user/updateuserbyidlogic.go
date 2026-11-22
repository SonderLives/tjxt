package userlogic

import (
	"context"
	"database/sql"
	"errors"

	"tjxt/apps/user/rpc/internal/model"
	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserByIdLogic {
	return &UpdateUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUserById 管理端更新指定用户：仅覆盖请求中的非零字段，避免误清零未传字段。
func (l *UpdateUserByIdLogic) UpdateUserById(in *pb.UserDTO) (*pb.EmptyResponse, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("用户不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新用户失败")
	}
	d, err := l.svcCtx.UserDetailModel.FindOne(l.ctx, in.Id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新用户失败")
	}
	if d == nil {
		d = &model.UserDetail{Id: in.Id, Type: u.Type, Creater: sql.NullInt64{Int64: in.Id, Valid: true}}
	}

	if in.Username != "" {
		u.Username = in.Username
	}
	if in.CellPhone != "" {
		u.CellPhone = in.CellPhone
	}
	if in.Type != 0 {
		u.Type = int64(in.Type)
	}
	u.Updater = in.Id

	if in.Name != "" {
		d.Name = in.Name
	}
	if in.Gender != 0 {
		d.Gender = int64(in.Gender)
	}
	d.Icon = applyStr(d.Icon, in.Icon)
	d.Email = applyStr(d.Email, in.Email)
	d.Qq = applyStr(d.Qq, in.Qq)
	d.Job = applyStr(d.Job, in.Job)
	d.Province = applyStr(d.Province, in.Province)
	d.City = applyStr(d.City, in.City)
	d.District = applyStr(d.District, in.District)
	d.Intro = applyStr(d.Intro, in.Intro)
	d.Photo = applyStr(d.Photo, in.Photo)
	if in.RoleId != 0 {
		d.RoleId = in.RoleId
	}
	d.Updater = in.Id

	if err := l.svcCtx.UserModel.Update(l.ctx, u); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新用户失败")
	}
	if err := l.svcCtx.UserDetailModel.Update(l.ctx, d); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新用户失败")
	}
	return &pb.EmptyResponse{}, nil
}
