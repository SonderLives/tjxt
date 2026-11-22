package medialogic

import (
	"fmt"
	"strings"

	"tjxt/common/idgen"
)

// mockBaseURL 对象存储未配置（Store==nil）时的 mock 占位地址。
// 配置 MinIO（S3 协议，兼容 COS/OSS 网关）后所有签名/地址均由 store 生成。
const mockBaseURL = "http://127.0.0.1:9000"

// objectKey 按文件名生成云端对象 key：media/<雪花ID><扩展名>。
func objectKey(filename string) string {
	return fmt.Sprintf("media/%d%s", idgen.NextID(), fileExt(filename))
}

// mockFileURL 生成文件的 mock 访问地址（Store==nil 回退用）。
func mockFileURL(key string) string {
	return mockBaseURL + "/" + key
}

// mockUploadURL 生成文件的 mock 上传地址（Store==nil 回退用）。
func mockUploadURL(key string) string {
	return mockBaseURL + "/upload/" + key
}

// mockPlayURL 生成媒资的 mock 播放地址（Store==nil 回退用）。
func mockPlayURL(idOrKey any) string {
	return fmt.Sprintf("%s/play/%v", mockBaseURL, idOrKey)
}

// supportedMediaTypes 支持的媒资类型白名单
var supportedMediaTypes = map[string]struct{}{
	"video": {},
	"image": {},
	"audio": {},
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
