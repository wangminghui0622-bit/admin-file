package logic

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"
	"github.com/wangminghui0622-bit/admin-file/internal/utils/filex"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UploadFile receives a client stream, writes the bytes to the public store
// and records the metadata in fms_files.
func (l *UploadFileLogic) UploadFile(stream file.File_UploadFileServer) error {
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

	typeName, typeNum := filex.SplitMIME(first.FileType)
	conf := l.svcCtx.Config.UploadConf

	dateStr := time.Now().Format("2006-01-02")
	publicDir := filepath.Join(conf.PublicStorePath, "fms", typeName, dateStr)
	privateDir := filepath.Join(conf.PrivateStorePath, "fms", typeName, dateStr)
	if err := os.MkdirAll(publicDir, 0o755); err != nil {
		logx.Errorw("failed to create public store directory", logx.Field("detail", err.Error()))
		return errorx.New(errorx.ServerError, "failed to create directory")
	}
	if err := os.MkdirAll(privateDir, 0o755); err != nil {
		logx.Errorw("failed to create private store directory", logx.Field("detail", err.Error()))
		return errorx.New(errorx.ServerError, "failed to create directory")
	}

	fileID := newUUID()
	storeName := fileID + "." + ext
	targetPath := filepath.Join(publicDir, storeName)

	target, err := os.Create(targetPath)
	if err != nil {
		logx.Errorw("failed to create file", logx.Field("detail", err.Error()))
		return errorx.New(errorx.ServerError, "failed to create file")
	}

	var total int64
	writeChunk := func(data []byte) error {
		if len(data) == 0 {
			return nil
		}
		total += int64(len(data))
		if err := filex.CheckOverSize(typeName, total, conf); err != nil {
			return err
		}
		if _, err := target.Write(data); err != nil {
			return errorx.New(errorx.ServerError, "failed to write file")
		}
		return nil
	}

	fail := func(e error) error {
		_ = target.Close()
		_ = os.Remove(targetPath)
		return e
	}

	if err := writeChunk(first.Data); err != nil {
		return fail(err)
	}
	if !first.Last {
		for {
			packet, recvErr := stream.Recv()
			if recvErr == io.EOF {
				break
			}
			if recvErr != nil {
				return fail(errorx.New(errorx.ServerError, recvErr.Error()))
			}
			if err := writeChunk(packet.Data); err != nil {
				return fail(err)
			}
			if packet.Last {
				break
			}
		}
	}
	if err := target.Close(); err != nil {
		return fail(errorx.New(errorx.ServerError, "failed to write file"))
	}
	if total == 0 {
		_ = os.Remove(targetPath)
		return errorx.New(errorx.InvalidArgument, "file is empty")
	}

	relativePath := fmt.Sprintf("/fms/%s/%s/%s", typeName, dateStr, storeName)
	userID := ""
	if first.UserId != nil {
		userID = *first.UserId
	}

	row := &model.FmsFile{
		ID:       fileID,
		Status:   ptr(int32(1)),
		Name:     baseName,
		FileType: int32(typeNum),
		Size:     total,
		Path:     relativePath,
		UserID:   userID,
		Md5:      first.Md5,
	}

	err = l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		q := query.Use(tx)
		if err := q.FmsFile.WithContext(l.ctx).Create(row); err != nil {
			return err
		}
		if first.TagId != 0 {
			if err := q.FileTagFile.WithContext(l.ctx).Create(&model.FileTagFile{
				FileTagID: int64(first.TagId),
				FileID:    fileID,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		_ = os.Remove(targetPath)
		return dbErr(err)
	}

	tagIDs := make([]uint64, 0, 1)
	if first.TagId != 0 {
		tagIDs = append(tagIDs, first.TagId)
	}
	resp := toPBFile(row, tagIDs, conf.ServerURL)
	// The reference upload response carries the original file name (with suffix).
	resp.Name = ptr(first.Name)

	return stream.SendAndClose(resp)
}
