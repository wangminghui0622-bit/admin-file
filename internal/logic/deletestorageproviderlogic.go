package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteStorageProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteStorageProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteStorageProviderLogic {
	return &DeleteStorageProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteStorageProviderLogic) DeleteStorageProvider(in *file.IDsReq) (*file.BaseResp, error) {
	if len(in.Ids) == 0 {
		return nil, errorx.New(errorx.InvalidArgument, "ids is required")
	}
	ids := make([]int64, 0, len(in.Ids))
	for _, id := range in.Ids {
		ids = append(ids, int64(id))
	}

	q := l.svcCtx.Query
	count, err := q.FmsCloudFile.WithContext(l.ctx).
		Where(q.FmsCloudFile.CloudFileStorageProviders.In(ids...)).Count()
	if err != nil {
		return nil, dbErr(err)
	}
	if count != 0 {
		return nil, errorx.New(errorx.InvalidArgument, "storage provider has files")
	}

	if _, err := q.FmsStorageProvider.WithContext(l.ctx).
		Where(q.FmsStorageProvider.ID.In(ids...)).Delete(); err != nil {
		return nil, dbErr(err)
	}

	l.svcCtx.CloudStorage.Reload()

	return successResp(), nil
}
