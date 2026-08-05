package searchlogic

import (
	"context"
	"strings"

	"tjxt/apps/search/rpc/internal/svc"
	"tjxt/apps/search/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveInterestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveInterestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveInterestsLogic {
	return &SaveInterestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveInterests 保存用户兴趣：校验格式（仅数字与逗号、去重、≤255），
// 不存在则插入、存在则更新（Upsert 一步完成）。
func (l *SaveInterestsLogic) SaveInterests(in *pb.SaveInterestsReq) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("用户id非法")
	}
	interests, err := normalizeInterests(in.Interests)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.InterestsModel.Upsert(l.ctx, in.Id, interests); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "保存兴趣失败")
	}
	return &pb.Empty{}, nil
}

// normalizeInterests 校验并规范化兴趣字符串：
// 仅允许数字与逗号，去重保持顺序，结果不超过 255 字符；空串表示清空。
func normalizeInterests(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}

	seen := make(map[string]struct{})
	kept := make([]string, 0, 8)
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if !isDigits(p) {
			return "", xerr.BadRequestf("兴趣格式非法，仅支持逗号分隔的数字分类id")
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		kept = append(kept, p)
	}

	out := strings.Join(kept, ",")
	if len(out) > 255 {
		return "", xerr.BadRequestf("兴趣分类id过长，最多255字符")
	}
	return out, nil
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
