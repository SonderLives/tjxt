package medialogic

import (
	"strconv"

	"tjxt/pkg/timex"
	"tjxt/rpc/media/internal/model"
	"tjxt/rpc/media/internal/store"
	"tjxt/rpc/media/pb"
)

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
		CreateTime: timex.Format(m.CreateTime),
		// UseTimes 表中无对应字段，真实场景由引用计数维护，此处统一返回 0
	}
}

// toFileVO 文件模型转 RPC VO
func toFileVO(m *model.File, st *store.Store) *pb.FileVO {
	path := mockFileURL(m.Key)
	if st != nil {
		path = st.PublicURL(m.Key)
	}
	return &pb.FileVO{
		Id:       m.Id,
		Key:      m.Key,
		Filename: m.Filename,
		Path:     path,
		Status:   int32(m.Status),
	}
}
