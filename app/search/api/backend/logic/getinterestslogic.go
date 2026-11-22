package searchlogic

import (
	"context"
	"errors"

	"tjxt/app/search/api/backend/model"
	"tjxt/app/search/api/backend/svc"
	"tjxt/app/search/api/backend/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInterestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInterestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInterestsLogic {
	return &GetInterestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetInterests 查询用户兴趣，未设置过时返回空 VO（不报错）。
func (l *GetInterestsLogic) GetInterests(in *pb.IdReq) (*pb.InterestsVO, error) {
	if in.Id <= 0 {
		return nil, errorx.BadRequestf("用户id非法")
	}

	data, err := l.svcCtx.InterestsModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &pb.InterestsVO{Id: in.Id}, nil
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询兴趣失败")
	}

	return &pb.InterestsVO{
		Id:         data.Id,
		Interests:  data.Interests.String,
		CreateTime: data.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: data.UpdateTime.Format("2006-01-02 15:04:05"),
	}, nil
}
