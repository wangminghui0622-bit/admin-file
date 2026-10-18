package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCloudFileTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCloudFileTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCloudFileTagLogic {
	return &UpdateCloudFileTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCloudFileTagLogic) UpdateCloudFileTag(in *file.CloudFileTagInfo) (*file.BaseResp, error) {
	if in.Id == nil {
		return nil, errorx.New(errorx.InvalidArgument, "id is required")
	}
	q := l.svcCtx.Query
	if _, err := q.FmsCloudFileTag.WithContext(l.ctx).Where(q.FmsCloudFileTag.ID.Eq(int64(*in.Id))).First(); err != nil {
		return nil, dbErr(err)
	}

	updates := map[string]any{}
	if in.Status != nil {
		updates["status"] = int32(*in.Status)
	}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Remark != nil {
		updates["remark"] = *in.Remark
	}
	if len(updates) > 0 {
		if _, err := q.FmsCloudFileTag.WithContext(l.ctx).
			Where(q.FmsCloudFileTag.ID.Eq(int64(*in.Id))).Updates(updates); err != nil {
			return nil, dbErr(err)
		}
	}
	return successResp(), nil
}
