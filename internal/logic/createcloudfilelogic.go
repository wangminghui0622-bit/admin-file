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

type CreateCloudFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCloudFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCloudFileLogic {
	return &CreateCloudFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCloudFileLogic) CreateCloudFile(in *file.CloudFileInfo) (*file.BaseUUIDResp, error) {
	if in.ProviderId == nil {
		return nil, errorx.New(errorx.InvalidArgument, "provider_id is required")
	}
	if in.Name == nil || *in.Name == "" {
		return nil, errorx.New(errorx.InvalidArgument, "name is required")
	}
	if in.Url == nil || *in.Url == "" {
		return nil, errorx.New(errorx.InvalidArgument, "url is required")
	}
	if in.UserId == nil || *in.UserId == "" {
		return nil, errorx.New(errorx.InvalidArgument, "user_id is required")
	}

	q := l.svcCtx.Query
	if _, err := q.FmsStorageProvider.WithContext(l.ctx).
		Where(q.FmsStorageProvider.ID.Eq(int64(*in.ProviderId))).First(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.New(errorx.InvalidArgument, "storage provider not exist")
		}
		return nil, dbErr(err)
	}

	state := true
	if in.State != nil {
		state = *in.State
	}
	var size int64
	if in.Size != nil {
		size = int64(*in.Size)
	}
	var fileType int32
	if in.FileType != nil {
		fileType = int32(*in.FileType)
	}
	providerID := int64(*in.ProviderId)
	fileID := newUUID()
	row := &model.FmsCloudFile{
		ID:                        fileID,
		State:                     &state,
		Name:                      *in.Name,
		URL:                       *in.Url,
		Size:                      size,
		FileType:                  fileType,
		UserID:                    *in.UserId,
		CloudFileStorageProviders: &providerID,
	}

	tagIDs := dedupeUint64(in.TagIds)
	err := l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		qtx := query.Use(tx)
		if err := qtx.FmsCloudFile.WithContext(l.ctx).Create(row); err != nil {
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

	return &file.BaseUUIDResp{Id: fileID, Msg: successMsg}, nil
}
