package logic

import (
	"context"
	"errors"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateCloudFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCloudFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCloudFileLogic {
	return &UpdateCloudFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCloudFileLogic) UpdateCloudFile(in *file.CloudFileInfo) (*file.BaseResp, error) {
	if in.Id == nil || *in.Id == "" {
		return nil, errorx.New(errorx.InvalidArgument, "id is required")
	}
	fileID := *in.Id
	q := l.svcCtx.Query
	if _, err := q.FmsCloudFile.WithContext(l.ctx).Where(q.FmsCloudFile.ID.Eq(fileID)).First(); err != nil {
		return nil, dbErr(err)
	}

	if in.ProviderId != nil {
		if _, err := q.FmsStorageProvider.WithContext(l.ctx).
			Where(q.FmsStorageProvider.ID.Eq(int64(*in.ProviderId))).First(); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errorx.New(errorx.InvalidArgument, "storage provider not exist")
			}
			return nil, dbErr(err)
		}
	}

	tagIDs := dedupeUint64(in.TagIds)
	err := l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		qtx := query.Use(tx)

		updates := map[string]any{}
		if in.State != nil {
			updates["state"] = *in.State
		}
		if in.Name != nil {
			updates["name"] = *in.Name
		}
		if in.Url != nil {
			updates["url"] = *in.Url
		}
		if in.Size != nil {
			updates["size"] = int64(*in.Size)
		}
		if in.FileType != nil {
			updates["file_type"] = int32(*in.FileType)
		}
		if in.UserId != nil {
			updates["user_id"] = *in.UserId
		}
		if in.ProviderId != nil {
			updates["cloud_file_storage_providers"] = int64(*in.ProviderId)
		}
		if len(updates) > 0 {
			if _, err := qtx.FmsCloudFile.WithContext(l.ctx).
				Where(qtx.FmsCloudFile.ID.Eq(fileID)).Updates(updates); err != nil {
				return err
			}
		}

		if _, err := qtx.CloudFileTagCloudFile.WithContext(l.ctx).
			Where(qtx.CloudFileTagCloudFile.CloudFileID.Eq(fileID)).Delete(); err != nil {
			return err
		}
		for _, tagID := range tagIDs {
			if err := qtx.CloudFileTagCloudFile.WithContext(l.ctx).Create(&model.CloudFileTagCloudFile{
				CloudFileTagID: int64(tagID),
				CloudFileID:    fileID,
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
