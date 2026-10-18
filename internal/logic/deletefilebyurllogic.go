package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"
	"github.com/wangminghui0622-bit/admin-file/internal/utils/filex"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileByUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileByUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileByUrlLogic {
	return &DeleteFileByUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileByUrlLogic) DeleteFileByUrl(in *file.DeleteByUrlReq) (*file.BaseResp, error) {
	id, err := filex.UUIDFromURL(in.Url)
	if err != nil {
		return nil, err
	}
	return NewDeleteFileLogic(l.ctx, l.svcCtx).DeleteFile(&file.UUIDsReq{Ids: []string{id}})
}
