package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStorageProviderByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStorageProviderByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStorageProviderByIdLogic {
	return &GetStorageProviderByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetStorageProviderByIdLogic) GetStorageProviderById(in *file.IDReq) (*file.StorageProviderInfo, error) {
	q := l.svcCtx.Query
	row, err := q.FmsStorageProvider.WithContext(l.ctx).Where(q.FmsStorageProvider.ID.Eq(int64(in.Id))).First()
	if err != nil {
		return nil, dbErr(err)
	}
	return toPBStorageProvider(row), nil
}
