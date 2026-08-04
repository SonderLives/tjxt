package logic

import (
	"context"

	"common/auth"
	"common/result"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 课程分类列表 ============

type CategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListLogic {
	return &CategoryListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CategoryListLogic) CategoryList(req *types.CategoryListQuery) (resp *result.R, err error) {
	list, err := l.svcCtx.CategoryService.List(l.ctx, req.Name, req.Status)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 新增课程分类 ============

type CategoryAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryAddLogic {
	return &CategoryAddLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CategoryAddLogic) CategoryAdd(req *types.CategoryAddDTO) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.CategoryService.Add(l.ctx, req, userID); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 获取课程分类信息 ============

type CategoryGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryGetLogic {
	return &CategoryGetLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CategoryGetLogic) CategoryGet(id int64) (resp *result.R, err error) {
	vo, err := l.svcCtx.CategoryService.Get(l.ctx, id)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}

// ============ 删除分类信息 ============

type CategoryDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryDeleteLogic {
	return &CategoryDeleteLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CategoryDeleteLogic) CategoryDelete(id int64) (resp *result.R, err error) {
	if err := l.svcCtx.CategoryService.Delete(l.ctx, id); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 课程分类停用或启用 ============

type CategoryDisableOrEnableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryDisableOrEnableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryDisableOrEnableLogic {
	return &CategoryDisableOrEnableLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CategoryDisableOrEnableLogic) CategoryDisableOrEnable(req *types.CategoryDisableOrEnableDTO) (resp *result.R, err error) {
	if err := l.svcCtx.CategoryService.DisableOrEnable(l.ctx, req); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 更新课程分类 ============

type CategoryUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryUpdateLogic {
	return &CategoryUpdateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CategoryUpdateLogic) CategoryUpdate(req *types.CategoryUpdateDTO) (resp *result.R, err error) {
	if err := l.svcCtx.CategoryService.Update(l.ctx, req); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 获取所有课程分类（分类树） ============

type CategoryAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryAllLogic {
	return &CategoryAllLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CategoryAllLogic) CategoryAll(admin bool) (resp *result.R, err error) {
	list, err := l.svcCtx.CategoryService.All(l.ctx, admin)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 获取所有课程分类（不分层） ============

type CategoryAllOfOneLevelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryAllOfOneLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryAllOfOneLevelLogic {
	return &CategoryAllOfOneLevelLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CategoryAllOfOneLevelLogic) CategoryAllOfOneLevel() (resp *result.R, err error) {
	list, err := l.svcCtx.CategoryService.AllOfOneLevel(l.ctx)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}
