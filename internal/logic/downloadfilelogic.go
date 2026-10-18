package logic

import (
	"context"
	"io"
	"os"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"
	"github.com/wangminghui0622-bit/admin-file/internal/utils/filex"

	"github.com/zeromicro/go-zero/core/logx"
)

const downloadChunkSize = 32 * 1024

type DownloadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDownloadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadFileLogic {
	return &DownloadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DownloadFile streams public files to anyone and private files to their owner.
func (l *DownloadFileLogic) DownloadFile(in *file.DownloadFileReq, stream file.File_DownloadFileServer) error {
	q := l.svcCtx.Query
	row, err := q.FmsFile.WithContext(l.ctx).Where(q.FmsFile.ID.Eq(in.Id)).First()
	if err != nil {
		return dbErr(err)
	}

	status := int32(1)
	if row.Status != nil {
		status = *row.Status
	}
	userID := ""
	if in.UserId != nil {
		userID = *in.UserId
	}

	var storePath string
	if status == 1 {
		storePath = l.svcCtx.Config.UploadConf.PublicStorePath
	} else {
		if row.UserID != userID {
			return errorx.New(errorx.Unauthorized, "无权访问此文件")
		}
		storePath = l.svcCtx.Config.UploadConf.PrivateStorePath
	}

	f, err := os.Open(filex.LocalPath(storePath, row.Path))
	if err != nil {
		logx.Errorw("failed to open file", logx.Field("detail", err.Error()), logx.Field("path", row.Path))
		return errorx.New(errorx.ServerError, "failed to open file")
	}
	defer f.Close()

	buf := make([]byte, downloadChunkSize)
	for {
		n, readErr := f.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			if err := stream.Send(&file.FileChunk{Data: chunk}); err != nil {
				return err
			}
		}
		if readErr == io.EOF {
			return nil
		}
		if readErr != nil {
			return errorx.New(errorx.ServerError, readErr.Error())
		}
	}
}
