package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileLogic {
	return &UpdateFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFileLogic) UpdateFile(in *file.UpdateFileReq) (*file.BaseResp, error) {
	if in.Id == "" {
		return nil, errorx.New(errorx.InvalidArgument, "id is required")
	}
	err := l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		q := query.Use(tx)
		if _, err := q.FmsFile.WithContext(l.ctx).Where(q.FmsFile.ID.Eq(in.Id)).First(); err != nil {
			return err
		}

		if in.Name != nil {
			if _, err := q.FmsFile.WithContext(l.ctx).
				Where(q.FmsFile.ID.Eq(in.Id)).Update(q.FmsFile.Name, *in.Name); err != nil {
				return err
			}
		}

		// repeated fields have no presence in proto3: always replace the tag set.
		if _, err := q.FileTagFile.WithContext(l.ctx).
			Where(q.FileTagFile.FileID.Eq(in.Id)).Delete(); err != nil {
			return err
		}
		for _, tagID := range dedupeUint64(in.TagIds) {
			if err := q.FileTagFile.WithContext(l.ctx).Create(&model.FileTagFile{
				FileTagID: int64(tagID),
				FileID:    in.Id,
			}); err != nil {
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
