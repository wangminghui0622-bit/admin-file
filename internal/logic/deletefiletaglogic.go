package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileTagLogic {
	return &DeleteFileTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileTagLogic) DeleteFileTag(in *file.IDsReq) (*file.BaseResp, error) {
	if len(in.Ids) == 0 {
		return nil, errorx.New(errorx.InvalidArgument, "ids is required")
	}
	ids := make([]int64, 0, len(in.Ids))
	for _, id := range in.Ids {
		ids = append(ids, int64(id))
	}
	q := l.svcCtx.Query
	if _, err := q.FmsFileTag.WithContext(l.ctx).Where(q.FmsFileTag.ID.In(ids...)).Delete(); err != nil {
		return nil, dbErr(err)
	}
	return successResp(), nil
}
