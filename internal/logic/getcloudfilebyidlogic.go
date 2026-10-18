package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudFileByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCloudFileByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFileByIdLogic {
	return &GetCloudFileByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCloudFileByIdLogic) GetCloudFileById(in *file.UUIDReq) (*file.CloudFileInfo, error) {
	q := l.svcCtx.Query
	row, err := q.FmsCloudFile.WithContext(l.ctx).Where(q.FmsCloudFile.ID.Eq(in.Id)).First()
	if err != nil {
		return nil, dbErr(err)
	}
	return toPBCloudFile(row, getCloudFileTagIDs(l.ctx, q, row.ID)), nil
}
