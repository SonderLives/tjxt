package service

import (
	"context"
	"database/sql"
	"sort"
	"strings"

	"tjxt/pkg/xerr"
	"tjxt/apps/course/api/internal/model"
	"tjxt/apps/course/api/internal/types"
)

// CategoryService 课程分类业务接口
type CategoryService interface {
	// List 查询课程分类信息（含课程数量、三级分类数量，可按名称/状态过滤）。
	List(ctx context.Context, name string, status int64) ([]*types.CategoryVO, error)
	// Add 新增课程分类。
	Add(ctx context.Context, req *types.CategoryAddDTO, userId int64) error
	// Get 获取课程分类信息。
	Get(ctx context.Context, id int64) (*types.CategoryInfoVO, error)
	// Delete 删除分类信息。
	Delete(ctx context.Context, id int64) error
	// DisableOrEnable 课程分类停用或启用（联动子分类）。
	DisableOrEnable(ctx context.Context, req *types.CategoryDisableOrEnableDTO) error
	// Update 更新课程分类。
	Update(ctx context.Context, req *types.CategoryUpdateDTO) error
	// All 获取所有课程分类信息（仅 id、名称、分类关系）。admin=false 仅返回启用且有课程的分类。
	All(ctx context.Context, admin bool) ([]*types.SimpleCategoryVO, error)
	// AllOfOneLevel 获取所有课程分类，不分层。
	AllOfOneLevel(ctx context.Context) ([]*types.CategoryVO, error)
}

type categoryService struct {
	categoryModel      model.CategoryModel
	courseModel        model.CourseModel
	courseDraftModel   model.CourseDraftModel
}

// NewCategoryService 创建课程分类业务服务。
func NewCategoryService(categoryModel model.CategoryModel, courseModel model.CourseModel, courseDraftModel model.CourseDraftModel) CategoryService {
	return &categoryService{categoryModel: categoryModel, courseModel: courseModel, courseDraftModel: courseDraftModel}
}

func (s *categoryService) List(ctx context.Context, name string, status int64) ([]*types.CategoryVO, error) {
	list, err := s.categoryModel.ListAll(ctx, true)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	if len(list) == 0 {
		return []*types.CategoryVO{}, nil
	}
	// 课程数量：正式 + 草稿（任一分类级别匹配计入）
	courseNumMap, err := s.countCourseNumOfCategory(ctx)
	if err != nil {
		return nil, err
	}
	thirdNumMap := statisticThirdCategory(list)
	vos := buildCategoryVOs(list, courseNumMap, thirdNumMap)
	if name == "" && status == 0 {
		return vos, nil
	}
	return filterCategoryVOs(vos, name, status), nil
}

// countCourseNumOfCategory 统计每个分类下的课程数量（正式+草稿）。
func (s *categoryService) countCourseNumOfCategory(ctx context.Context) (map[int64]int64, error) {
	result := make(map[int64]int64)
	courses, err := s.courseModel.ListAll(ctx)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "统计课程分类课程数量失败")
	}
	for _, c := range courses {
		addCount(result, c.FirstCateId, c.SecondCateId, c.ThirdCateId)
	}
	drafts, err := s.courseDraftModel.ListAll(ctx)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "统计课程分类课程数量失败")
	}
	for _, c := range drafts {
		addCount(result, c.FirstCateId, c.SecondCateId, c.ThirdCateId)
	}
	return result, nil
}

func addCount(m map[int64]int64, cates ...int64) {
	for _, id := range cates {
		if id > 0 {
			m[id]++
		}
	}
}

// statisticThirdCategory 统计每个一级/二级分类拥有的三级分类数量。
func statisticThirdCategory(list []*model.Category) map[int64]int64 {
	result := make(map[int64]int64)
	// 三级分类按父（二级）分组
	secondOfThird := make(map[int64]int64) // 二级 id -> 三级数量
	thirdParent := make(map[int64]int64)   // 三级 id -> 二级 id
	for _, c := range list {
		if c.Level == 3 {
			secondOfThird[c.ParentId]++
			thirdParent[c.Id] = c.ParentId
		}
	}
	for _, c := range list {
		if c.Level == 2 {
			result[c.Id] = secondOfThird[c.Id]
		}
	}
	// 一级分类的三级数量 = 其下所有二级分类的三级数量之和
	for _, c := range list {
		if c.Level == 3 {
			second := thirdParent[c.Id]
			for _, c2 := range list {
				if c2.Level == 2 && c2.Id == second {
					result[c2.ParentId]++
				}
			}
		}
	}
	return result
}

func buildCategoryVOs(list []*model.Category, courseNum, thirdNum map[int64]int64) []*types.CategoryVO {
	// id -> VO
	nodeMap := make(map[int64]*types.CategoryVO, len(list))
	for _, c := range list {
		vo := &types.CategoryVO{
			Id:        c.Id,
			Name:      c.Name,
			ParentId:  c.ParentId,
			Level:     int32(c.Level),
			Index:     int32(c.Priority),
			Status:    int32(c.Status),
			StatusDesc: statusDesc(c.Status),
			CreateTime: c.CreateTime.Format(timeFormat),
			UpdateTime: c.UpdateTime.Format(timeFormat),
			CourseNum: int32(courseNum[c.Id]),
			ThirdCategoryNum: int32(thirdNum[c.Id]),
		}
		nodeMap[c.Id] = vo
	}
	roots := make([]*types.CategoryVO, 0)
	for _, c := range list {
		vo := nodeMap[c.Id]
		if c.ParentId == CategoryRoot {
			roots = append(roots, vo)
		} else if parent, ok := nodeMap[c.ParentId]; ok {
			parent.Children = append(parent.Children, vo)
		} else {
			roots = append(roots, vo)
		}
	}
	return roots
}

// filterCategoryVOs 按名称关键字和状态过滤分类树：当前分类通过 或 存在通过的子分类 则保留。
func filterCategoryVOs(list []*types.CategoryVO, name string, status int64) []*types.CategoryVO {
	result := make([]*types.CategoryVO, 0, len(list))
	for _, vo := range list {
		if filterCategoryVO(vo, name, status) {
			result = append(result, vo)
		}
	}
	return result
}

func filterCategoryVO(vo *types.CategoryVO, name string, status int64) bool {
	pass := true
	if status != 0 {
		pass = int64(vo.Status) == status
	}
	if pass && name != "" {
		pass = vo.Name != "" && strings.Contains(vo.Name, name)
	}
	if !pass && len(vo.Children) == 0 {
		return false
	}
	children := make([]*types.CategoryVO, 0, len(vo.Children))
	for _, child := range vo.Children {
		if filterCategoryVO(child, name, status) {
			children = append(children, child)
		}
	}
	vo.Children = children
	return pass || len(vo.Children) > 0
}

func (s *categoryService) Add(ctx context.Context, req *types.CategoryAddDTO, userId int64) error {
	if req.Name == "" {
		return xerr.BadRequestf("分类名称不能为空")
	}
	// 校验同父分类下是否有同名分类
	if err := s.checkSameName(ctx, req.ParentId, req.Name, 0); err != nil {
		return err
	}
	level := int64(1)
	if req.ParentId != CategoryRoot {
		parent, err := s.categoryModel.FindById(ctx, req.ParentId)
		if err == sql.ErrNoRows {
			return xerr.Conflict("父分类不存在")
		}
		if err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "查询父分类失败")
		}
		if parent.Level >= 3 {
			return xerr.Conflict("三级分类下不能再创建子分类")
		}
		level = parent.Level + 1
	}
	if level > 3 {
		return xerr.Conflict("分类级别最多三级")
	}
	_, err := s.categoryModel.Insert(ctx, &model.Category{
		Name:   req.Name,
		ParentId: req.ParentId,
		Level:  level,
		Priority: int64(req.Index),
		Status: CategoryDisable, // 新增分类默认禁用
		Creater: userId,
		Updater: userId,
	})
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "新增课程分类失败")
	}
	return nil
}

func (s *categoryService) Get(ctx context.Context, id int64) (*types.CategoryInfoVO, error) {
	category, err := s.categoryModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return &types.CategoryInfoVO{}, nil
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	vo := &types.CategoryInfoVO{
		Id:            category.Id,
		Name:          category.Name,
		Index:         int32(category.Priority),
		CategoryLevel: int32(category.Level),
		Status:        int32(category.Status),
		StatusDesc:    statusDesc(category.Status),
		CreateTime:    category.CreateTime.Format(timeFormat),
		UpdateTime:    category.UpdateTime.Format(timeFormat),
	}
	firstCategoryId := int64(0)
	if category.Level == 3 {
		second, err := s.categoryModel.FindById(ctx, category.ParentId)
		if err == nil {
			vo.SecondCategoryName = second.Name
			firstCategoryId = second.ParentId
		}
	} else if category.Level == 2 {
		firstCategoryId = category.ParentId
	}
	if firstCategoryId != 0 {
		first, err := s.categoryModel.FindById(ctx, firstCategoryId)
		if err == nil {
			vo.FirstCategoryName = first.Name
		}
	}
	return vo, nil
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	// 有子分类不能删除
	children, err := s.categoryModel.ListByParent(ctx, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询子分类失败")
	}
	if len(children) > 0 {
		return xerr.Conflict("该分类下存在子分类，无法删除")
	}
	category, err := s.categoryModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return xerr.NotFound("分类不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询分类失败")
	}
	// 分类下有课程不能删除
	courseNum, err := s.courseModel.CountByCategoryId(ctx, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "统计分类课程数量失败")
	}
	if courseNum > 0 {
		return xerr.Conflict("该分类下存在课程，无法删除")
	}
	if err := s.categoryModel.DeleteById(ctx, id); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除分类失败")
	}
	_ = category
	return nil
}

func (s *categoryService) DisableOrEnable(ctx context.Context, req *types.CategoryDisableOrEnableDTO) error {
	category, err := s.categoryModel.FindById(ctx, req.Id)
	if err == sql.ErrNoRows {
		return xerr.Conflict("课程分类不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	status := int64(req.Status)
	// 启用时父分类必须已启用
	if status == CategoryEnable && category.ParentId != CategoryRoot {
		parent, err := s.categoryModel.FindById(ctx, category.ParentId)
		if err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "查询父分类失败")
		}
		if parent.Status == CategoryDisable {
			return xerr.Conflict("父分类禁用中，无法启用当前分类")
		}
	}
	// 联动子分类：一级禁用禁用直接子分类和孙子分类，二级禁用禁用直接子分类
	childIds := make([]int64, 0)
	directChildren, err := s.categoryModel.ListByParent(ctx, req.Id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询子分类失败")
	}
	for _, c := range directChildren {
		childIds = append(childIds, c.Id)
	}
	grandChildren := make([]int64, 0)
	for _, c := range directChildren {
		gc, err := s.categoryModel.ListByParent(ctx, c.Id)
		if err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "查询子分类失败")
		}
		for _, g := range gc {
			grandChildren = append(grandChildren, g.Id)
		}
	}
	childIds = append(childIds, grandChildren...)

	// 更新自身与联动分类
	if err := s.categoryModel.UpdateById(ctx, &model.Category{Id: req.Id, Status: status}, "status"); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新分类状态失败")
	}
	if len(childIds) > 0 {
		if err := s.categoryModel.UpdateStatusByIDs(ctx, childIds, status); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "更新子分类状态失败")
		}
	}
	return nil
}

func (s *categoryService) Update(ctx context.Context, req *types.CategoryUpdateDTO) error {
	category, err := s.categoryModel.FindById(ctx, req.Id)
	if err == sql.ErrNoRows {
		return xerr.NotFound("分类不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询分类失败")
	}
	if err := s.checkSameName(ctx, category.ParentId, req.Name, req.Id); err != nil {
		return err
	}
	if err := s.categoryModel.UpdateById(ctx, &model.Category{
		Id:       req.Id,
		Name:     req.Name,
		Priority: int64(req.Index),
	}, "name", "priority"); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新分类失败")
	}
	return nil
}

// checkSameName 校验同一父分类下（或与父分类同名）是否存在同名分类。
func (s *categoryService) checkSameName(ctx context.Context, parentId int64, name string, currentId int64) error {
	num, err := s.categoryModel.CountByNameSameSibling(ctx, parentId, name)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "校验分类名称失败")
	}
	if num > 0 {
		if currentId > 0 {
			// 更新时同名但为自身则允许
			self, err := s.categoryModel.FindById(ctx, currentId)
			if err != nil || self.Name != name {
				return xerr.Conflict("分类名称已存在")
			}
			return nil
		}
		return xerr.Conflict("分类名称已存在")
	}
	return nil
}

func (s *categoryService) All(ctx context.Context, admin bool) ([]*types.SimpleCategoryVO, error) {
	var list []*model.Category
	var err error
	if admin {
		list, err = s.categoryModel.ListAll(ctx, true)
	} else {
		list, err = s.categoryModel.ListEnabled(ctx)
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	if len(list) == 0 {
		return []*types.SimpleCategoryVO{}, nil
	}
	// 非 admin 仅返回有课程的分类
	withCourse := map[int64]bool{}
	if !admin {
		courseNumMap, err := s.countCourseNumOfCategory(ctx)
		if err != nil {
			return nil, err
		}
		for id := range courseNumMap {
			withCourse[id] = true
		}
		filtered := make([]*model.Category, 0, len(list))
		for _, c := range list {
			if withCourse[c.Id] {
				filtered = append(filtered, c)
			}
		}
		list = filtered
	}
	roots := buildSimpleCategoryVOs(list)
	// 过滤掉没有三级分类的分支
	return filterSimpleCategories(roots), nil
}

func buildSimpleCategoryVOs(list []*model.Category) []*types.SimpleCategoryVO {
	nodeMap := make(map[int64]*types.SimpleCategoryVO, len(list))
	for _, c := range list {
		nodeMap[c.Id] = &types.SimpleCategoryVO{
			Id:       c.Id,
			Name:     c.Name,
			Level:    int32(c.Level),
			ParentId: c.ParentId,
		}
	}
	roots := make([]*types.SimpleCategoryVO, 0)
	for _, c := range list {
		vo := nodeMap[c.Id]
		if c.ParentId == CategoryRoot {
			roots = append(roots, vo)
		} else if parent, ok := nodeMap[c.ParentId]; ok {
			parent.Children = append(parent.Children, vo)
		} else {
			roots = append(roots, vo)
		}
	}
	sortSimple(roots)
	return roots
}

func sortSimple(list []*types.SimpleCategoryVO) {
	sort.SliceStable(list, func(i, j int) bool { return list[i].Id > list[j].Id })
	for _, v := range list {
		sortSimple(v.Children)
	}
}

// filterSimpleCategories 移除没有三级子分类的分支。
func filterSimpleCategories(list []*types.SimpleCategoryVO) []*types.SimpleCategoryVO {
	result := make([]*types.SimpleCategoryVO, 0, len(list))
	for _, vo := range list {
		if vo.Level == 3 {
			result = append(result, vo)
			continue
		}
		children := filterSimpleCategories(vo.Children)
		if len(children) > 0 {
			vo.Children = children
			result = append(result, vo)
		}
	}
	return result
}

func (s *categoryService) AllOfOneLevel(ctx context.Context) ([]*types.CategoryVO, error) {
	list, err := s.categoryModel.ListAll(ctx, true)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	thirdNum := statisticThirdCategory(list)
	result := make([]*types.CategoryVO, 0, len(list))
	for _, c := range list {
		result = append(result, &types.CategoryVO{
			Id:              c.Id,
			Name:            c.Name,
			ParentId:        c.ParentId,
			Level:           int32(c.Level),
			Index:           int32(c.Priority),
			Status:          int32(c.Status),
			StatusDesc:      statusDesc(c.Status),
			ThirdCategoryNum: int32(thirdNum[c.Id]),
			CreateTime:      c.CreateTime.Format(timeFormat),
			UpdateTime:      c.UpdateTime.Format(timeFormat),
		})
	}
	return result, nil
}

func statusDesc(status int64) string {
	if status == CategoryEnable {
		return "启用"
	}
	return "禁用"
}
