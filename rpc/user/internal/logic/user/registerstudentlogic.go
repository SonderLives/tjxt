package userlogic

import (
	"context"
	"database/sql"

	"tjxt/apps/user/rpc/internal/model"
	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"
	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterStudentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterStudentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterStudentLogic {
	return &RegisterStudentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RegisterStudent 学员注册：校验 -> 查重 -> 写 user + user_detail（密码 bcrypt 加密）。
func (l *RegisterStudentLogic) RegisterStudent(in *pb.StudentFormRequest) (*pb.IdResponse, error) {
	if in.CellPhone == "" || in.Password == "" {
		return nil, xerr.BadRequestf("手机号与密码均必填")
	}
	exists, err := l.svcCtx.UserModel.ExistsByCellPhone(l.ctx, in.CellPhone)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "注册失败")
	}
	if exists {
		return nil, xerr.Conflict("该手机号已注册")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "注册失败")
	}
	id := idgen.NextID()
	user := &model.User{
		Id:        id,
		Username:  in.CellPhone,
		CellPhone: in.CellPhone,
		Password:  string(hash),
		Type:      2,
		Status:    1,
		Creater:   sql.NullInt64{Int64: id, Valid: true},
		Updater:   id,
	}
	if _, err := l.svcCtx.UserModel.Insert(l.ctx, user); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "注册失败")
	}
	detail := &model.UserDetail{
		Id:      id,
		Type:    2,
		Creater: sql.NullInt64{Int64: id, Valid: true},
		Updater: id,
	}
	if _, err := l.svcCtx.UserDetailModel.Insert(l.ctx, detail); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "注册失败")
	}
	return &pb.IdResponse{Id: id}, nil
}
