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
	"golang.org/x/crypto/bcrypt"
)

type UpdateCurrentUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCurrentUserLogic {
	return &UpdateCurrentUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateCurrentUser 更新当前登录用户：合并非零资料字段；若携带新密码则校验原密码后重设。
func (l *UpdateCurrentUserLogic) UpdateCurrentUser(in *pb.UserFormRequest) (*pb.EmptyResponse, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("用户不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新失败")
	}
	d, err := l.svcCtx.UserDetailModel.FindOne(l.ctx, in.Id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新失败")
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

	if in.Password != "" {
		if in.OldPassword == "" {
			return nil, xerr.BadRequestf("修改密码需提供原密码")
		}
		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.OldPassword)) != nil {
			return nil, xerr.Unauthorized("原密码错误")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "更新失败")
		}
		u.Password = string(hash)
	}

	if err := l.svcCtx.UserModel.Update(l.ctx, u); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新失败")
	}
	if err := l.svcCtx.UserDetailModel.Update(l.ctx, d); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新失败")
	}
	return &pb.EmptyResponse{}, nil
}
