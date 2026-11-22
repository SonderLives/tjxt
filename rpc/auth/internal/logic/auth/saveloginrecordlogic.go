package authlogic

import (
	"context"
	"database/sql"
	"time"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveLoginRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveLoginRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveLoginRecordLogic {
	return &SaveLoginRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveLoginRecord 落库一条登录记录。
// 该调用属于登录链路的旁路统计，不应阻断主流程，因此写库失败仅记录日志。
func (l *SaveLoginRecordLogic) SaveLoginRecord(in *pb.LoginRecordReq) (*pb.Empty, error) {
	if in.UserId <= 0 {
		return nil, xerr.BadRequestf("用户 id 无效")
	}

	now := time.Now()
	record := &model.LoginRecord{
		UserId:    in.UserId,
		CellPhone: sql.NullString{String: in.CellPhone, Valid: in.CellPhone != ""},
		LoginTime: now,
		LoginDate: sql.NullTime{Time: now, Valid: true},
		Ipv4:      in.Ipv4,
	}
	if _, err := l.svcCtx.LoginRecordModel.Insert(l.ctx, record); err != nil {
		l.Errorf("save login record failed, userId=%d: %v", in.UserId, err)
	}
	return &pb.Empty{}, nil
}
