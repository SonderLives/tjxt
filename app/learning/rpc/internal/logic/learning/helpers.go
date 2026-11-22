package learninglogic

import (
	"context"

	"tjxt/common/auth"
)

// lessonRequestUserID 解析本次请求的学习者 id：
// 请求字段优先（供 api/course 等上游以字段显式指定学员），
// 未携带时回退到 gRPC metadata 透传的登录用户 id（上游 API 经
// auth.OutgoingUserIdInterceptor 注入，见 pkg/auth/rpc.go）。
func lessonRequestUserID(ctx context.Context, requestUserID int64) (int64, error) {
	if requestUserID > 0 {
		return requestUserID, nil
	}
	return auth.UserIdFromIncomingMD(ctx)
}
