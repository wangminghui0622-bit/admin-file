package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudFileListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCloudFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFileListLogic {
	return &GetCloudFileListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCloudFileListLogic) GetCloudFileList(in *file.CloudFileListReq) (*file.CloudFileListResp, error) {
	q := l.svcCtx.Query
	query := q.FmsCloudFile.WithContext(l.ctx)

	if in.Name != nil {
		query = query.Where(q.FmsCloudFile.Name.Like("%" + *in.Name + "%"))
	}
	if in.ProviderId != nil {
		query = query.Where(q.FmsCloudFile.CloudFileStorageProviders.Eq(int64(*in.ProviderId)))
	}
	if in.ProviderName != nil {
		providers, err := q.FmsStorageProvider.WithContext(l.ctx).
			Select(q.FmsStorageProvider.ID).
			Where(q.FmsStorageProvider.Name.Eq(*in.ProviderName)).Find()
		if err != nil {
			return nil, dbErr(err)
		}
		if len(providers) == 0 {
			return &file.CloudFileListResp{Data: []*file.CloudFileInfo{}}, nil
		}
		ids := make([]int64, 0, len(providers))
		for _, provider := range providers {
			ids = append(ids, provider.ID)
		}
		query = query.Where(q.FmsCloudFile.CloudFileStorageProviders.In(ids...))
	}
	if len(in.TagIds) > 0 {
		ids, err := cloudFileIDsByTagIDs(l.ctx, q, dedupeUint64(in.TagIds))
		if err != nil {
			return nil, dbErr(err)
		}
		if len(ids) == 0 {
			return &file.CloudFileListResp{Data: []*file.CloudFileInfo{}}, nil
		}
		query = query.Where(q.FmsCloudFile.ID.In(ids...))
	}
	if in.FileType != nil && *in.FileType != 0 {
		query = query.Where(q.FmsCloudFile.FileType.Eq(int32(*in.FileType)))
	}

	total, err := query.Count()
	if err != nil {
		return nil, dbErr(err)
	}
	offset, limit := paging(in.Page, in.PageSize)
	rows, err := query.Order(q.FmsCloudFile.ID.Desc()).Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, dbErr(err)
	}

	fileIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		fileIDs = append(fileIDs, row.ID)
	}
	tagMap := getCloudFileTagIDsBatch(l.ctx, q, fileIDs)

	resp := &file.CloudFileListResp{Total: uint64(total), Data: make([]*file.CloudFileInfo, 0, len(rows))}
	for _, row := range rows {
		resp.Data = append(resp.Data, toPBCloudFile(row, tagMap[row.ID]))
	}
	return resp, nil
}
