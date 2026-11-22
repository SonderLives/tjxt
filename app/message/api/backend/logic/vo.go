package messagelogic

import (
	"tjxt/service/message/api/backend/model"
	"tjxt/service/message/api/backend/pb"
	"tjxt/common/timex"
)

func toNoticeTemplateVO(m *model.NoticeTemplate) *pb.NoticeTemplateVO {
	return &pb.NoticeTemplateVO{
		Id:            m.Id,
		Name:          m.Name,
		Code:          m.Code,
		Type:          int32(m.Type),
		Status:        int32(m.Status),
		Title:         m.Title.String,
		Content:       m.Content,
		IsSmsTemplate: m.IsSmsTemplate == 1,
		CreateTime:    timex.Format(m.CreateTime),
	}
}

func toNoticeTaskVO(m *model.NoticeTask) *pb.NoticeTaskVO {
	return &pb.NoticeTaskVO{
		Id:         m.Id,
		TemplateId: m.TemplateId,
		Name:       m.Name,
		Partial:    m.Partial == 1,
		PushTime:   timex.FormatNullTime(m.PushTime),
		Interval:   int32(m.Interval.Int64),
		ExpireTime: timex.FormatNullTime(m.ExpireTime),
		MaxTimes:   int32(m.MaxTimes),
		Finished:   m.Finished == 1,
		CreateTime: timex.Format(m.CreateTime),
	}
}

func toMessageTemplateVO(m *model.MessageTemplate) *pb.MessageTemplateVO {
	return &pb.MessageTemplateVO{
		Id:                m.Id,
		Name:              m.Name,
		PlatformCode:      m.PlatformCode,
		SignName:          m.SignName,
		ThirdTemplateCode: m.ThirdTemplateCode,
		Content:           m.Content,
		TemplateId:        m.TemplateId,
		Status:            int32(m.Status),
		CreateTime:        timex.Format(m.CreateTime),
	}
}

func toUserInboxVO(m *model.UserInbox) *pb.UserInboxVO {
	return &pb.UserInboxVO{
		Id:         m.Id,
		UserId:     m.UserId,
		Type:       int32(m.Type),
		Title:      m.Title,
		Content:    m.Content,
		IsRead:     m.IsRead == 1,
		Publisher:  m.Publisher,
		PushTime:   timex.Format(m.PushTime),
		ExpireTime: timex.Format(m.ExpireTime),
	}
}
