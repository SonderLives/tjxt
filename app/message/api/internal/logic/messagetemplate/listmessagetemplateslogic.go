// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package messagetemplate

import (
	"context"

	backendlogic "tjxt/service/message/api/backend/logic"
	"tjxt/service/message/api/backend/pb"
	"tjxt/service/message/api/internal/svc"
	"tjxt/service/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMessageTemplatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMessageTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessageTemplatesLogic {
	return &ListMessageTemplatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMessageTemplatesLogic) ListMessageTemplates(req *types.PageReq) (resp *types.MessageTemplateListVO, err error) {
	r, err := backendlogic.NewListMessageTemplatesLogic(l.ctx, l.svcCtx.Backend).ListMessageTemplates(&pb.PageReq{
		PageNo:   int32(req.PageNo),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.MessageTemplateVO, 0, len(r.List))
	for _, item := range r.List {
		list = append(list, types.MessageTemplateVO{
			Id:                item.Id,
			Name:              item.Name,
			PlatformCode:      item.PlatformCode,
			SignName:          item.SignName,
			ThirdTemplateCode: item.ThirdTemplateCode,
			Content:           item.Content,
			TemplateId:        item.TemplateId,
			Status:            item.Status,
			CreateTime:        item.CreateTime,
		})
	}
	return &types.MessageTemplateListVO{Total: r.Total, List: list}, nil
}
