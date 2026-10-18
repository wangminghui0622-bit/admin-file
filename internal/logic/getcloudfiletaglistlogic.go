package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudFileTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCloudFileTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFileTagListLogic {
	return &GetCloudFileTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCloudFileTagListLogic) GetCloudFileTagList(in *file.CloudFileTagListReq) (*file.CloudFileTagListResp, error) {
	q := l.svcCtx.Query
	query := q.FmsCloudFileTag.WithContext(l.ctx)
	if in.Name != nil {
		query = query.Where(q.FmsCloudFileTag.Name.Like("%" + *in.Name + "%"))
	}
	if in.Remark != nil {
		query = query.Where(q.FmsCloudFileTag.Remark.Like("%" + *in.Remark + "%"))
	}

	total, err := query.Count()
	if err != nil {
		return nil, dbErr(err)
	}
	offset, limit := paging(in.Page, in.PageSize)
	rows, err := query.Order(q.FmsCloudFileTag.ID.Desc()).Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, dbErr(err)
	}

	resp := &file.CloudFileTagListResp{Total: uint64(total), Data: make([]*file.CloudFileTagInfo, 0, len(rows))}
	for _, row := range rows {
		resp.Data = append(resp.Data, toPBCloudFileTag(row))
	}
	return resp, nil
}
