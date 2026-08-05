package userlogic

import (
	"context"
	"database/sql"
	"errors"

	"tjxt/apps/user/rpc/internal/model"
	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"
	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type AddUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddUser 管理端新增用户：查重 -> 写 user + user_detail。
// 密码不随请求传入，统一使用默认初始口令（defaultInitialPassword），后续由用户自行修改或管理员重置。
func (l *AddUserLogic) AddUser(in *pb.UserDTO) (*pb.IdResponse, error) {
	if in.CellPhone == "" && in.Username == "" {
		return nil, xerr.BadRequestf("手机号与用户名不能都为空")
	}
	if in.CellPhone != "" {
		exists, err := l.svcCtx.UserModel.ExistsByCellPhone(l.ctx, in.CellPhone)
		if err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "创建用户失败")
		}
		if exists {
			return nil, xerr.Conflict("该手机号已存在")
		}
	}
	if in.Username != "" {
		if _, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username); err == nil {
			return nil, xerr.Conflict("该用户名已存在")
		} else if !errors.Is(err, model.ErrNotFound) {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "创建用户失败")
		}
	}

	userType := int64(in.Type)
	if userType == 0 {
		userType = 2
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultInitialPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "创建用户失败")
	}
	id := idgen.NextID()
	user := &model.User{
		Id:        id,
		Username:  in.Username,
		CellPhone: in.CellPhone,
		Password:  string(hash),
		Type:      userType,
		Status:    1,
		Creater:   sql.NullInt64{Int64: id, Valid: true},
		Updater:   id,
	}
	if _, err := l.svcCtx.UserModel.Insert(l.ctx, user); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "创建用户失败")
	}
	detail := &model.UserDetail{
		Id:        id,
		Type:      userType,
		Name:      in.Name,
		Gender:    int64(in.Gender),
		Icon:      nullStr(in.Icon),
		Email:     nullStr(in.Email),
		Qq:        nullStr(in.Qq),
		Job:       nullStr(in.Job),
		Province:  nullStr(in.Province),
		City:      nullStr(in.City),
		District:  nullStr(in.District),
		Intro:     nullStr(in.Intro),
		Photo:     nullStr(in.Photo),
		RoleId:    in.RoleId,
		Creater:   sql.NullInt64{Int64: id, Valid: true},
		Updater:   id,
	}
	if _, err := l.svcCtx.UserDetailModel.Insert(l.ctx, detail); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "创建用户失败")
	}
	return &pb.IdResponse{Id: id}, nil
}
