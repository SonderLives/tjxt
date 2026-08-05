package logic

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tjxt/apps/media/rpc/internal/model"
	"tjxt/apps/media/rpc/pb"
	"tjxt/pkg/utils/idgen"
)

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

// 支持的媒资类型
var supportedMediaTypes = map[string]struct{}{
	"video": {},
	"image": {},
	"audio": {},
}

// mockBaseURL 对象存储的 mock 占位地址。
// 真实项目中替换为腾讯云 COS / 阿里 OSS 的访问域名（如 https://xxx.cos.ap-guangzhou.myqcloud.com）。
const mockBaseURL = "http://127.0.0.1:9000"

// mockToken 生成 mock 签名串。
// 真实项目中由 COS/OSS 的 STS 临时密钥服务签发，包含 ak/sk/token/有效期等。
func mockToken() string {
	sum := md5.Sum([]byte(strconv.FormatInt(idgen.NextID(), 10)))
	return hex.EncodeToString(sum[:])
}

// mockObjectKey 按文件名生成云端 key：mock/<雪花ID><文件扩展名>。
// 真实项目中 key 由 COS/OSS 上传策略生成（目录 + 雪花ID + 扩展名）。
func mockObjectKey(filename string) string {
	return fmt.Sprintf("mock/%d%s", idgen.NextID(), fileExt(filename))
}

// mockFileURL 生成文件的 mock 访问地址。
// 真实项目返回 COS/OSS 的持久化访问地址。
func mockFileURL(key string) string {
	return mockBaseURL + "/" + key
}

// mockUploadURL 生成文件的 mock 上传地址。
// 真实项目返回 COS/OSS 的 PUT 直传地址或预签名上传 URL。
func mockUploadURL(key string) string {
	return mockBaseURL + "/upload/" + key
}

// mockPlayURL 生成媒资的 mock 播放地址。
func mockPlayURL(idOrKey any) string {
	return fmt.Sprintf("%s/play/%v", mockBaseURL, idOrKey)
}

// supportedMediaType 校验媒资类型是否在白名单（video/image/audio）内
func supportedMediaType(t string) bool {
	_, ok := supportedMediaTypes[t]
	return ok
}

// fileExt 提取文件名扩展名（含点，如 ".mp4"），无扩展名返回空串
func fileExt(name string) string {
	if i := strings.LastIndex(name, "."); i > 0 && i < len(name)-1 {
		return name[i:]
	}
	return ""
}

// formatTime 时间格式化为 "2006-01-02 15:04:05"
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// isNotFound 判断错误是否为数据不存在
func isNotFound(err error) bool {
	return err == sql.ErrNoRows || err == model.ErrNotFound
}

// toMediaVO 媒资模型转 RPC VO
func toMediaVO(m *model.Media) *pb.MediaVO {
	return &pb.MediaVO{
		Id:         m.Id,
		Filename:   m.Filename,
		MediaUrl:   m.MediaUrl,
		CoverUrl:   m.CoverUrl,
		Duration:   m.Duration,
		Size:       m.Size,
		Status:     int32(m.Status),
		Creater:    strconv.FormatInt(m.Creater, 10),
		CreateTime: formatTime(m.CreateTime),
		// UseTimes 表中无对应字段，真实场景由引用计数维护，此处统一返回 0
	}
}

// toFileVO 文件模型转 RPC VO
func toFileVO(m *model.File) *pb.FileVO {
	return &pb.FileVO{
		Id:       m.Id,
		Key:      m.Key,
		Filename: m.Filename,
		// Path 表中无对应字段，真实场景为 OSS/COS 的访问地址，此处 mock 生成
		Path:   mockFileURL(m.Key),
		Status: int32(m.Status),
	}
}
