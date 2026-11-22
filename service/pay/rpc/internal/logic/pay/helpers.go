package paylogic

import (
	"database/sql"

	"tjxt/common/page"
	"tjxt/service/pay/rpc/internal/model"
	"tjxt/service/pay/rpc/pb"
)

// formatNullString 取 sql.NullString 有效值，无效返回空串。
func formatNullString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// normalizePage 把 proto 的分页参数归一化成 offset/limit
func normalizePage(p *pb.PageQueryPayChannelsRequest) (offset, limit int64) {
	return page.Normalize(p.PageNo, p.PageSize)
}

// buildPayChannel 把 proto 请求转换为 model 数据结构（供新增/更新复核用）
func buildPayChannel(in *pb.PayChannelRequest) *model.PayChannel {
	return &model.PayChannel{
		Id:              in.Id,
		Name:            in.Name,
		ChannelCode:     in.ChannelCode,
		ChannelPriority: int64(in.ChannelPriority),
		ChannelIcon:     in.ChannelIcon,
		Status:          PayChannelStatusEnabled,
	}
}
