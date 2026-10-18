package logic

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteCloudFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCloudFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCloudFileLogic {
	return &DeleteCloudFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCloudFileLogic) DeleteCloudFile(in *file.UUIDsReq) (*file.BaseResp, error) {
	if len(in.Ids) == 0 {
		return nil, errorx.New(errorx.InvalidArgument, "ids is required")
	}

	q := l.svcCtx.Query
	rows, err := q.FmsCloudFile.WithContext(l.ctx).Where(q.FmsCloudFile.ID.In(in.Ids...)).Find()
	if err != nil {
		return nil, dbErr(err)
	}

	if l.svcCtx.Config.UploadConf.DeleteFileWithCloud {
		for _, row := range rows {
			if row.CloudFileStorageProviders == nil {
				continue
			}
			provider, err := q.FmsStorageProvider.WithContext(l.ctx).
				Where(q.FmsStorageProvider.ID.Eq(*row.CloudFileStorageProviders)).First()
			if err != nil {
				logx.Errorw("failed to find the provider of the cloud file", logx.Field("fileId", row.ID))
				continue
			}
			client, pd, ok := l.svcCtx.CloudStorage.Client(provider.Name)
			if !ok {
				logx.Errorw("the provider client is not loaded", logx.Field("provider", provider.Name))
				continue
			}

			var key string
			if !pd.UseCdn {
				parts := strings.Split(row.URL, pd.Endpoint)
				if len(parts) != 2 {
					logx.Errorw("failed to find the key of the cloud file", logx.Field("url", row.URL))
					return nil, errorx.New(errorx.ServerError, "failed to find the key of the cloud file")
				}
				key = parts[1]
			} else {
				key = strings.TrimPrefix(row.URL, pd.CdnURL)
				if len(key) < 2 {
					logx.Errorw("failed to find the key of the cloud file", logx.Field("url", row.URL))
					return nil, errorx.New(errorx.ServerError, "failed to find the key of the cloud file")
				}
			}

			if _, err := client.DeleteObjectWithContext(l.ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(pd.Bucket),
				Key:    aws.String(key),
			}); err != nil {
				logx.Errorw("failed to delete the cloud file", logx.Field("detail", err.Error()))
				return nil, errorx.New(errorx.ServerError, "failed to delete the cloud file")
			}
		}
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		qtx := query.Use(tx)
		if _, err := qtx.CloudFileTagCloudFile.WithContext(l.ctx).
			Where(qtx.CloudFileTagCloudFile.CloudFileID.In(in.Ids...)).Delete(); err != nil {
			return err
		}
		if _, err := qtx.FmsCloudFile.WithContext(l.ctx).
			Where(qtx.FmsCloudFile.ID.In(in.Ids...)).Delete(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, dbErr(err)
	}
	return successResp(), nil
}
