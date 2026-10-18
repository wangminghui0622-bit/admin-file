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

type SetFileStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetFileStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetFileStatusLogic {
	return &SetFileStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetFileStatus moves a file between the public and private stores and updates
// its status atomically; a failed move rolls the DB update back.
func (l *SetFileStatusLogic) SetFileStatus(in *file.StatusReq) (*file.BaseResp, error) {
	if in.Status != 1 && in.Status != 2 {
		return nil, errorx.New(errorx.InvalidArgument, "invalid status")
	}
	conf := l.svcCtx.Config.UploadConf

	err := l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		q := query.Use(tx)
		row, err := q.FmsFile.WithContext(l.ctx).Where(q.FmsFile.ID.Eq(in.Id)).First()
		if err != nil {
			return err
		}

		current := int32(1)
		if row.Status != nil {
			current = *row.Status
		}
		if current == int32(in.Status) {
			return nil
		}

		if _, err := q.FmsFile.WithContext(l.ctx).
			Where(q.FmsFile.ID.Eq(in.Id)).Update(q.FmsFile.Status, int32(in.Status)); err != nil {
			return err
		}

		var src, dst string
		if in.Status == 1 {
			src = filex.LocalPath(conf.PrivateStorePath, row.Path)
			dst = filex.LocalPath(conf.PublicStorePath, row.Path)
		} else {
			src = filex.LocalPath(conf.PublicStorePath, row.Path)
			dst = filex.LocalPath(conf.PrivateStorePath, row.Path)
		}
		if err := os.Rename(src, dst); err != nil {
			logx.Errorw("failed to change the path of file", logx.Field("detail", err.Error()))
			return err
		}
		return nil
	})
	if err != nil {
		return nil, dbErr(err)
	}
	return successResp(), nil
}
