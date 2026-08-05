// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package inbox

import (
	"context"
	"strconv"

	"tjxt/pkg/auth"

	"google.golang.org/grpc/metadata"
)

// metadataKeyUserID 与 message-rpc 侧的 metadataKeyUserID 保持一致
const metadataKeyUserID = "user_id"

// ctxWithUserId 将当前登录用户 id 附加到 RPC 出站 metadata，
// 供 message-rpc 做站内信归属校验（MarkInboxRead/DeleteInbox）。
func ctxWithUserId(ctx context.Context) (context.Context, error) {
	userId, err := auth.UserIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return metadata.AppendToOutgoingContext(ctx, metadataKeyUserID, strconv.FormatInt(userId, 10)), nil
}
