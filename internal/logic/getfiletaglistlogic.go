package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileTagListLogic {
	return &GetFileTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileTagListLogic) GetFileTagList(in *file.FileTagListReq) (*file.FileTagListResp, error) {
	q := l.svcCtx.Query
	query := q.FmsFileTag.WithContext(l.ctx)
	if in.Name != nil {
		query = query.Where(q.FmsFileTag.Name.Like("%" + *in.Name + "%"))
	}
	if in.Remark != nil {
		query = query.Where(q.FmsFileTag.Remark.Like("%" + *in.Remark + "%"))
	}

	total, err := query.Count()
	if err != nil {
		return nil, dbErr(err)
	}
	offset, limit := paging(in.Page, in.PageSize)
	rows, err := query.Order(q.FmsFileTag.ID.Desc()).Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, dbErr(err)
	}

	resp := &file.FileTagListResp{Total: uint64(total), Data: make([]*file.FileTagInfo, 0, len(rows))}
	for _, row := range rows {
		resp.Data = append(resp.Data, toPBFileTag(row))
	}
	return resp, nil
}
