package service

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"common/xerr"
	"course/internal/model"
	"course/internal/types"
)

// CatalogueService 章节目录业务接口（正式表只读接口）。
type CatalogueService interface {
	// BatchQuery 根据章节目录批量查询基础信息。
	BatchQuery(ctx context.Context, ids []int64) ([]*types.CataSimpleInfoVO, error)
	// QuerySectionInfoById 获取小节信息。
	QuerySectionInfoById(ctx context.Context, id int64) (*types.CataSimpleInfoVO, error)
	// GetSimpleSectionInfo 获取小节对应的媒资 id 与课程 id（内部调用）。
	GetSimpleSectionInfo(ctx context.Context, sectionId int64) (*types.SectionInfoDTO, error)
	// CountMediaUserInfo 统计媒资被引用次数（内部调用）。
	CountMediaUserInfo(ctx context.Context, mediaIds []int64) ([]*types.MediaQuoteDTO, error)
	// GetCatasIndexList 根据课程 id 查询所有章节的序号。
	GetCatasIndexList(ctx context.Context, courseId int64) ([]*types.CataSimpleInfoVO, error)
}

type catalogueService struct {
	cataModel       *model.CourseCatalogueModel
	trailerDuration int64
}

// NewCatalogueService 创建章节目录业务服务。
func NewCatalogueService(cataModel *model.CourseCatalogueModel, trailerDuration int64) CatalogueService {
	return &catalogueService{cataModel: cataModel, trailerDuration: trailerDuration}
}

// BatchQuery 根据章节目录批量查询基础信息。
func (s *catalogueService) BatchQuery(ctx context.Context, ids []int64) ([]*types.CataSimpleInfoVO, error) {
	if len(ids) == 0 {
		return []*types.CataSimpleInfoVO{}, nil
	}
	catas, err := s.cataModel.ListByIds(ctx, ids)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询章节目录失败")
	}
	result := make([]*types.CataSimpleInfoVO, 0, len(catas))
	for _, c := range catas {
		result = append(result, &types.CataSimpleInfoVO{
			Id:     c.Id,
			Name:   c.Name,
			CIndex: int32(c.CIndex),
		})
	}
	return result, nil
}

// QuerySectionInfoById 获取小节信息。
func (s *catalogueService) QuerySectionInfoById(ctx context.Context, id int64) (*types.CataSimpleInfoVO, error) {
	cata, err := s.cataModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return &types.CataSimpleInfoVO{}, nil
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询小节失败")
	}
	if cata.Type != CataTypeSection {
		return &types.CataSimpleInfoVO{}, nil
	}
	vo := &types.CataSimpleInfoVO{
		Id:     cata.Id,
		Name:   cata.Name,
		CIndex: int32(cata.CIndex),
	}
	// 查询章信息
	chapter, err := s.cataModel.FindById(ctx, cata.ParentCatalogueId)
	if err == nil {
		vo.ChapterIndex = int32(chapter.CIndex)
	}
	return vo, nil
}

// GetSimpleSectionInfo 获取小节对应的媒资 id 与课程 id。
func (s *catalogueService) GetSimpleSectionInfo(ctx context.Context, sectionId int64) (*types.SectionInfoDTO, error) {
	if sectionId <= 0 {
		return nil, xerr.BadRequestf("小节 id 不能为空")
	}
	cata, err := s.cataModel.FindById(ctx, sectionId)
	if err == sql.ErrNoRows {
		return &types.SectionInfoDTO{}, nil
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询小节失败")
	}
	if cata.Type != CataTypeSection {
		return &types.SectionInfoDTO{}, nil
	}
	freeDuration := int64(0)
	if cata.Trailer == 1 {
		freeDuration = s.trailerDuration
	}
	return &types.SectionInfoDTO{
		CourseId:     cata.CourseId,
		MediaId:      cata.MediaId,
		Trailer:      cata.Trailer == 1,
		FreeDuration: int32(freeDuration),
	}, nil
}

// CountMediaUserInfo 统计媒资被引用次数。
func (s *catalogueService) CountMediaUserInfo(ctx context.Context, mediaIds []int64) ([]*types.MediaQuoteDTO, error) {
	result := make([]*types.MediaQuoteDTO, 0, len(mediaIds))
	quoteMap := make(map[int64]int64, len(mediaIds))
	if len(mediaIds) > 0 {
		catas, err := s.cataModel.ListByMediaIds(ctx, mediaIds)
		if err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询媒资引用失败")
		}
		for _, c := range catas {
			quoteMap[c.MediaId]++
		}
	}
	for _, mediaId := range mediaIds {
		result = append(result, &types.MediaQuoteDTO{MediaId: mediaId, QuoteNum: int32(quoteMap[mediaId])})
	}
	return result, nil
}

// GetCatasIndexList 根据课程 id 查询所有章节的序号。
func (s *catalogueService) GetCatasIndexList(ctx context.Context, courseId int64) ([]*types.CataSimpleInfoVO, error) {
	catas, err := s.cataModel.ListByCourseId(ctx, courseId, false)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询章节目录失败")
	}
	if len(catas) == 0 {
		return []*types.CataSimpleInfoVO{}, nil
	}
	// 章 id -> 章序号
	chapterMap := make(map[int64]int64, len(catas))
	for _, c := range catas {
		if c.Type == CataTypeChapter {
			chapterMap[c.Id] = c.CIndex
		}
	}
	result := make([]*types.CataSimpleInfoVO, 0)
	for _, c := range catas {
		if c.Type == CataTypeChapter {
			continue
		}
		result = append(result, &types.CataSimpleInfoVO{
			Id:     c.Id,
			Name:   c.Name,
			Index:  fmt.Sprintf("%d-%d", chapterMap[c.ParentCatalogueId], c.CIndex),
			CIndex: int32(c.CIndex),
		})
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
}
