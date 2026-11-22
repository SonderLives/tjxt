package userlogic

import (
	"context"
	"errors"

	"tjxt/app/user/rpc/internal/model"
	"tjxt/app/user/rpc/internal/svc"
	"tjxt/app/user/rpc/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UpdateStudentPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStudentPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStudentPasswordLogic {
	return &UpdateStudentPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateStudentPassword 学员修改密码（按手机号 + 类型=学员定位）。
func (l *UpdateStudentPasswordLogic) UpdateStudentPassword(in *pb.StudentFormRequest) (*pb.EmptyResponse, error) {
	if in.CellPhone == "" || in.Password == "" {
		return nil, errorx.BadRequestf("手机号与密码均必填")
	}
	user, err := l.svcCtx.UserModel.FindOneByCellPhoneType(l.ctx, in.CellPhone, 2)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("用户不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "修改密码失败")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "修改密码失败")
	}
	user.Password = string(hash)
	user.Updater = user.Id
	if err := l.svcCtx.UserModel.Update(l.ctx, user); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "修改密码失败")
	}
	return &pb.EmptyResponse{}, nil
}
