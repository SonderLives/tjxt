package medialogic

import "time"

// 文件状态
const (
	FileStatusPending  = 1 // 待上传
	FileStatusUploaded = 2 // 已上传未使用
	FileStatusUsed     = 3 // 已使用
)

// 媒资状态
const (
	MediaStatusUploading = 1 // 上传中
	MediaStatusUploaded  = 2 // 已上传
)

// 文件平台
const (
	PlatformTencent = 1 // 腾讯云 COS
	PlatformAliyun  = 2 // 阿里云 OSS
)

// PresignUploadTTL 预签名上传（PUT 直传）URL 有效期。
const PresignUploadTTL = 30 * time.Minute

// PresignViewTTL 预签名播放/预览（GET）URL 有效期（带时效防盗链）。
const PresignViewTTL = 2 * time.Hour
