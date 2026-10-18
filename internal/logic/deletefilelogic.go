package logic

import (
	"context"
	"os"

	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"
	"github.com/wangminghui0622-bit/admin-file/internal/utils/filex"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *file.UUIDsReq) (*file.BaseResp, error) {
	if len(in.Ids) == 0 {
		return nil, errorx.New(errorx.InvalidArgument, "ids is required")
	}
	err := l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		q := query.Use(tx)
		rows, err := q.FmsFile.WithContext(l.ctx).Where(q.FmsFile.ID.In(in.Ids...)).Find()
		if err != nil {
			return err
		}

		if _, err := q.FileTagFile.WithContext(l.ctx).
			Where(q.FileTagFile.FileID.In(in.Ids...)).Delete(); err != nil {
			return err
		}
		if _, err := q.FmsFile.WithContext(l.ctx).
			Where(q.FmsFile.ID.In(in.Ids...)).Delete(); err != nil {
			return err
		}

		conf := l.svcCtx.Config.UploadConf
		for _, row := range rows {
			storePath := conf.PublicStorePath
			if row.Status != nil && *row.Status != 1 {
				storePath = conf.PrivateStorePath
			}
			if err := os.RemoveAll(filex.LocalPath(storePath, row.Path)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, dbErr(err)
	}
	return successResp(), nil
}
