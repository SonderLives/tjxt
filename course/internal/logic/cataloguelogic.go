package logic

import (
	"context"

	"common/result"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 根据章节目录批量查询基础信息 ============

type CatalogueBatchQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCatalogueBatchQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CatalogueBatchQueryLogic {
	return &CatalogueBatchQueryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CatalogueBatchQueryLogic) BatchQuery(req *types.IdsQuery) (resp *result.R, err error) {
	list, err := l.svcCtx.CatalogueService.BatchQuery(l.ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 获取小节信息 ============

type CatalogueSectionInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCatalogueSectionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CatalogueSectionInfoLogic {
	return &CatalogueSectionInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CatalogueSectionInfoLogic) SectionInfo(id int64) (resp *result.R, err error) {
	vo, err := l.svcCtx.CatalogueService.QuerySectionInfoById(l.ctx, id)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}
