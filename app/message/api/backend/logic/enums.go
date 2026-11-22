package messagelogic

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
