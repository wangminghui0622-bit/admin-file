package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStorageProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStorageProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStorageProviderLogic {
	return &UpdateStorageProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateStorageProviderLogic) UpdateStorageProvider(in *file.StorageProviderInfo) (*file.BaseResp, error) {
	if in.Id == nil {
		return nil, errorx.New(errorx.InvalidArgument, "id is required")
	}
	q := l.svcCtx.Query
	id := int64(*in.Id)
	if _, err := q.FmsStorageProvider.WithContext(l.ctx).Where(q.FmsStorageProvider.ID.Eq(id)).First(); err != nil {
		return nil, dbErr(err)
	}

	updates := map[string]any{}
	if in.State != nil {
		updates["state"] = *in.State
	}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Bucket != nil {
		updates["bucket"] = *in.Bucket
	}
	if in.SecretId != nil {
		updates["secret_id"] = *in.SecretId
	}
	if in.SecretKey != nil {
		updates["secret_key"] = *in.SecretKey
	}
	if in.Endpoint != nil {
		updates["endpoint"] = *in.Endpoint
	}
	if in.Folder != nil {
		updates["folder"] = *in.Folder
	}
	if in.Region != nil {
		updates["region"] = *in.Region
	}
	if in.IsDefault != nil {
		updates["is_default"] = *in.IsDefault
	}
	if in.UseCdn != nil {
		updates["use_cdn"] = *in.UseCdn
	}
	if in.CdnUrl != nil {
		updates["cdn_url"] = *in.CdnUrl
	}
	if len(updates) > 0 {
		if _, err := q.FmsStorageProvider.WithContext(l.ctx).
			Where(q.FmsStorageProvider.ID.Eq(id)).Updates(updates); err != nil {
			return nil, dbErr(err)
		}
	}

	l.svcCtx.CloudStorage.Reload()

	return successResp(), nil
}
