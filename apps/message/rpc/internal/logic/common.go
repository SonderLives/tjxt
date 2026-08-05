package logic

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"tjxt/apps/message/rpc/internal/model"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/utils/idgen"

	"google.golang.org/grpc/metadata"
)

// metadataKeyUserID 上游 API 通过 gRPC metadata 透传的当前登录用户 id
const metadataKeyUserID = "user_id"

// timeLayout 时间统一格式 yyyy-MM-dd HH:mm:ss
const timeLayout = "2006-01-02 15:04:05"

// 通知模板状态（与表注释一致）
const (
	NoticeTemplateStatusDraft   = 0 // 草稿
	NoticeTemplateStatusInUse   = 1 // 使用中
	NoticeTemplateStatusStopped = 2 // 停用
)

// 短信模板状态（与表注释一致）
const (
	MessageTemplateStatusDisabled = 0 // 禁用
	MessageTemplateStatusEnabled  = 1 // 启用
)

// isNotFound 判断错误是否是数据不存在
func isNotFound(err error) bool {
	return err == sql.ErrNoRows || err == model.ErrNotFound
}

// userIdFromCtx 从 gRPC 入站 metadata 中读取上游透传的登录用户 id，缺失时返回 0
func userIdFromCtx(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}
	vals := md.Get(metadataKeyUserID)
	if len(vals) == 0 {
		return 0
	}
	id, err := strconv.ParseInt(vals[0], 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// nextID 雪花算法生成分布式 ID
func nextID() int64 { return idgen.NextID() }

// parseDateTime 解析 yyyy-MM-dd HH:mm:ss 格式时间，空串返回 NULL
func parseDateTime(s string) (sql.NullTime, error) {
	if s == "" {
		return sql.NullTime{}, nil
	}
	t, err := time.ParseInLocation(timeLayout, s, time.Local)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(timeLayout)
}

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
		CreateTime:    formatTime(m.CreateTime),
	}
}

func toNoticeTaskVO(m *model.NoticeTask) *pb.NoticeTaskVO {
	return &pb.NoticeTaskVO{
		Id:         m.Id,
		TemplateId: m.TemplateId,
		Name:       m.Name,
		Partial:    m.Partial == 1,
		PushTime:   formatNullTime(m.PushTime),
		Interval:   int32(m.Interval.Int64),
		ExpireTime: formatNullTime(m.ExpireTime),
		MaxTimes:   int32(m.MaxTimes),
		Finished:   m.Finished == 1,
		CreateTime: formatTime(m.CreateTime),
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
		CreateTime:        formatTime(m.CreateTime),
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
		PushTime:   formatTime(m.PushTime),
		ExpireTime: formatTime(m.ExpireTime),
	}
}
