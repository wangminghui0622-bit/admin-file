package logic

import (
	"context"
	"time"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gen/field"
)

type FileListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FileListLogic) FileList(in *file.FileListReq) (*file.FileListResp, error) {
	q := l.svcCtx.Query
	query := q.FmsFile.WithContext(l.ctx)

	// Privacy filter: public files for everyone, private files only for the owner.
	me := ""
	if in.UserId != nil {
		me = *in.UserId
	}
	query = query.Where(field.Or(
		q.FmsFile.Status.Eq(int32(1)),
		field.And(q.FmsFile.Status.Eq(int32(2)), q.FmsFile.UserID.Eq(me)),
	))

	if in.FileType != nil && *in.FileType != 0 {
		query = query.Where(q.FmsFile.FileType.Eq(int32(*in.FileType)))
	}
	if in.FileName != nil {
		query = query.Where(q.FmsFile.Name.Like("%" + *in.FileName + "%"))
	}
	if len(in.TagIds) > 0 {
		ids, err := fileIDsByTagIDs(l.ctx, q, dedupeUint64(in.TagIds))
		if err != nil {
			return nil, dbErr(err)
		}
		if len(ids) == 0 {
			return &file.FileListResp{Data: []*file.FileInfo{}}, nil
		}
		query = query.Where(q.FmsFile.ID.In(ids...))
	}
	if in.Status != nil {
		query = query.Where(q.FmsFile.Status.Eq(int32(*in.Status)))
	}
	if in.StartAt != nil {
		query = query.Where(q.FmsFile.CreatedAt.Gte(time.UnixMilli(*in.StartAt)))
	}
	if in.EndAt != nil {
		query = query.Where(q.FmsFile.CreatedAt.Lte(time.UnixMilli(*in.EndAt)))
	}

	total, err := query.Count()
	if err != nil {
		return nil, dbErr(err)
	}
	offset, limit := paging(in.Page, in.PageSize)
	rows, err := query.Order(q.FmsFile.ID.Desc()).Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, dbErr(err)
	}

	fileIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		fileIDs = append(fileIDs, row.ID)
	}
	tagMap := getFileTagIDsBatch(l.ctx, q, fileIDs)

	resp := &file.FileListResp{Total: uint64(total), Data: make([]*file.FileInfo, 0, len(rows))}
	for _, row := range rows {
		resp.Data = append(resp.Data, toPBFile(row, tagMap[row.ID], l.svcCtx.Config.UploadConf.ServerURL))
	}
	return resp, nil
}
