package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStorageProviderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStorageProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStorageProviderListLogic {
	return &GetStorageProviderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetStorageProviderListLogic) GetStorageProviderList(in *file.StorageProviderListReq) (*file.StorageProviderListResp, error) {
	q := l.svcCtx.Query
	query := q.FmsStorageProvider.WithContext(l.ctx)
	if in.Name != nil {
		query = query.Where(q.FmsStorageProvider.Name.Like("%" + *in.Name + "%"))
	}

	total, err := query.Count()
	if err != nil {
		return nil, dbErr(err)
	}
	offset, limit := paging(in.Page, in.PageSize)
	rows, err := query.Order(q.FmsStorageProvider.ID.Desc()).Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, dbErr(err)
	}

	resp := &file.StorageProviderListResp{Total: uint64(total), Data: make([]*file.StorageProviderInfo, 0, len(rows))}
	for _, row := range rows {
		resp.Data = append(resp.Data, toPBStorageProvider(row))
	}
	return resp, nil
}
