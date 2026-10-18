package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateStorageProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateStorageProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStorageProviderLogic {
	return &CreateStorageProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateStorageProviderLogic) CreateStorageProvider(in *file.StorageProviderInfo) (*file.BaseIDResp, error) {
	required := []struct {
		name  string
		value *string
	}{
		{"name", in.Name},
		{"bucket", in.Bucket},
		{"secret_id", in.SecretId},
		{"secret_key", in.SecretKey},
		{"endpoint", in.Endpoint},
		{"region", in.Region},
	}
	for _, item := range required {
		if item.value == nil || *item.value == "" {
			return nil, errorx.New(errorx.InvalidArgument, item.name+" is required")
		}
	}

	state := true
	if in.State != nil {
		state = *in.State
	}
	row := &model.FmsStorageProvider{
		State:     &state,
		Name:      *in.Name,
		Bucket:    *in.Bucket,
		SecretID:  *in.SecretId,
		SecretKey: *in.SecretKey,
		Endpoint:  *in.Endpoint,
		Folder:    in.Folder,
		Region:    *in.Region,
		IsDefault: in.IsDefault != nil && *in.IsDefault,
		UseCdn:    in.UseCdn != nil && *in.UseCdn,
		CdnURL:    in.CdnUrl,
	}
	if err := l.svcCtx.Query.FmsStorageProvider.WithContext(l.ctx).Create(row); err != nil {
		return nil, dbErr(err)
	}

	l.svcCtx.CloudStorage.Reload()

	return &file.BaseIDResp{Id: uint64(row.ID), Msg: successMsg}, nil
}
