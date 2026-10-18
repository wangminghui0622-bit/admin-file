package logic

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/cloud"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"
	"github.com/wangminghui0622-bit/admin-file/internal/utils/filex"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UploadCloudFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadCloudFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCloudFileLogic {
	return &UploadCloudFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UploadCloudFile spools the stream to a temp file (the S3 SDK needs an
// io.ReadSeeker), uploads it to the selected provider and records metadata.
func (l *UploadCloudFileLogic) UploadCloudFile(stream file.File_UploadCloudFileServer) error {
	first, err := stream.Recv()
	if err != nil {
		return errorx.New(errorx.InvalidArgument, "failed to receive file data")
	}
	if first.Name == "" {
		return errorx.New(errorx.InvalidArgument, "file name is required")
	}
	baseName, ext, ok := filex.SplitName(first.Name)
	if !ok {
		return errorx.New(errorx.InvalidArgument, "wrong file type")
	}

	var providerData cloud.ProviderData
	if first.Provider != "" {
		pd, exist := l.svcCtx.CloudStorage.Provider(first.Provider)
		if !exist {
			// 参考实现：显式指定但未注册的提供商静默回退默认提供商
			first.Provider = ""
		} else {
			providerData = pd
		}
	}
	if first.Provider == "" {
		pd, exist := l.svcCtx.CloudStorage.Default()
		if !exist {
			return errorx.New(errorx.InvalidArgument, "no default storage provider")
		}
		providerData = pd
	}
	client, providerData, ok := l.svcCtx.CloudStorage.Client(providerData.Name)
	if !ok {
		return errorx.New(errorx.InvalidArgument, "storage provider not exist")
	}

	typeName, typeNum := filex.SplitMIME(first.FileType)

	tmp, err := os.CreateTemp("", "admin-file-cloud-*")
	if err != nil {
		return errorx.New(errorx.ServerError, "failed to create temp file")
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()

	var total int64
	writeChunk := func(data []byte) error {
		if len(data) == 0 {
			return nil
		}
		total += int64(len(data))
		if err := filex.CheckOverSize(typeName, total, l.svcCtx.Config.UploadConf); err != nil {
			return err
		}
		if _, err := tmp.Write(data); err != nil {
			return errorx.New(errorx.ServerError, "failed to write file")
		}
		return nil
	}

	if err := writeChunk(first.Data); err != nil {
		return err
	}
	if !first.Last {
		for {
			packet, recvErr := stream.Recv()
			if recvErr == io.EOF {
				break
			}
			if recvErr != nil {
				return errorx.New(errorx.ServerError, recvErr.Error())
			}
			if err := writeChunk(packet.Data); err != nil {
				return err
			}
			if packet.Last {
				break
			}
		}
	}
	if total == 0 {
		return errorx.New(errorx.InvalidArgument, "file is empty")
	}
	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return errorx.New(errorx.ServerError, "failed to read file")
	}

	fileID := newUUID()
	storeName := fileID + "." + ext
	dateStr := time.Now().Format("2006-01-02")
	relativeSrc := fmt.Sprintf("%s/%s/%s/%s", providerData.Folder, dateStr, typeName, storeName)

	if _, err := client.PutObjectWithContext(l.ctx, &s3.PutObjectInput{
		Bucket: aws.String(providerData.Bucket),
		Key:    aws.String(relativeSrc),
		Body:   tmp,
	}); err != nil {
		logx.Errorw("failed to upload object", logx.Field("detail", err.Error()))
		return errorx.New(errorx.ServerError, "failed to upload object")
	}

	url := fmt.Sprintf("https://%s.%s%s", providerData.Bucket, providerData.Endpoint, relativeSrc)
	if providerData.UseCdn && providerData.CdnURL != "" {
		url = providerData.CdnURL + relativeSrc
	}

	userID := ""
	if first.UserId != nil {
		userID = *first.UserId
	}
	state := true
	providerID := int64(providerData.ID)
	row := &model.FmsCloudFile{
		ID:                        fileID,
		State:                     &state,
		Name:                      baseName,
		URL:                       url,
		Size:                      total,
		FileType:                  int32(typeNum),
		UserID:                    userID,
		CloudFileStorageProviders: &providerID,
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		q := query.Use(tx)
		if err := q.FmsCloudFile.WithContext(l.ctx).Create(row); err != nil {
			return err
		}
		if first.TagId != 0 {
			if err := q.CloudFileTagCloudFile.WithContext(l.ctx).Create(&model.CloudFileTagCloudFile{
				CloudFileTagID: int64(first.TagId),
				CloudFileID:    fileID,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return dbErr(err)
	}

	tagIDs := make([]uint64, 0, 1)
	if first.TagId != 0 {
		tagIDs = append(tagIDs, first.TagId)
	}
	resp := toPBCloudFile(row, tagIDs)
	resp.RelativeSrc = ptr(relativeSrc)

	return stream.SendAndClose(resp)
}
