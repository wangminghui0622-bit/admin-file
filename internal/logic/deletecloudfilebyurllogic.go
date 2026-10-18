package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"
	"github.com/wangminghui0622-bit/admin-file/internal/utils/filex"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCloudFileByUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCloudFileByUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCloudFileByUrlLogic {
	return &DeleteCloudFileByUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCloudFileByUrlLogic) DeleteCloudFileByUrl(in *file.DeleteByUrlReq) (*file.BaseResp, error) {
	id, err := filex.UUIDFromURL(in.Url)
	if err != nil {
		return nil, err
	}
	return NewDeleteCloudFileLogic(l.ctx, l.svcCtx).DeleteCloudFile(&file.UUIDsReq{Ids: []string{id}})
}
