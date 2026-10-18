package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileTagByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileTagByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileTagByIdLogic {
	return &GetFileTagByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileTagByIdLogic) GetFileTagById(in *file.IDReq) (*file.FileTagInfo, error) {
	q := l.svcCtx.Query
	row, err := q.FmsFileTag.WithContext(l.ctx).Where(q.FmsFileTag.ID.Eq(int64(in.Id))).First()
	if err != nil {
		return nil, dbErr(err)
	}
	return toPBFileTag(row), nil
}
