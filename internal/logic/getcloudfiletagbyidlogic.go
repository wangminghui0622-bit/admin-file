package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudFileTagByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCloudFileTagByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFileTagByIdLogic {
	return &GetCloudFileTagByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCloudFileTagByIdLogic) GetCloudFileTagById(in *file.IDReq) (*file.CloudFileTagInfo, error) {
	q := l.svcCtx.Query
	row, err := q.FmsCloudFileTag.WithContext(l.ctx).Where(q.FmsCloudFileTag.ID.Eq(int64(in.Id))).First()
	if err != nil {
		return nil, dbErr(err)
	}
	return toPBCloudFileTag(row), nil
}
