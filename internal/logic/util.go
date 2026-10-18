package logic

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/dal/query"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"gorm.io/gorm"
)

const (
	successMsg      = "success"
	defaultPageSize = 10
)

func successResp() *file.BaseResp {
	return &file.BaseResp{Msg: successMsg}
}

// dbErr maps gorm errors to the contract error codes.
func dbErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorx.New(errorx.NotExist, "record not found")
	}
	return errorx.New(errorx.ServerError, err.Error())
}

func paging(page, pageSize uint64) (offset, limit int) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return int((page - 1) * pageSize), int(pageSize)
}

func newUUID() string {
	return uuid.NewString()
}

func ptr[T any](v T) *T {
	return &v
}

func dedupeUint64(ids []uint64) []uint64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[uint64]struct{}, len(ids))
	result := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func toPBFile(row *model.FmsFile, tagIDs []uint64, serverURL string) *file.FileInfo {
	status := uint32(1)
	if row.Status != nil {
		status = uint32(*row.Status)
	}
	publicPath := ""
	if status == 1 {
		publicPath = serverURL + row.Path
	}
	return &file.FileInfo{
		Id:         ptr(row.ID),
		CreatedAt:  ptr(row.CreatedAt.UnixMilli()),
		UpdatedAt:  ptr(row.UpdatedAt.UnixMilli()),
		Status:     ptr(status),
		Name:       ptr(row.Name),
		FileType:   ptr(uint32(row.FileType)),
		Size:       ptr(uint64(row.Size)),
		Path:       ptr(row.Path),
		UserId:     ptr(row.UserID),
		PublicPath: ptr(publicPath),
		TagIds:     tagIDs,
	}
}

func toPBFileTag(row *model.FmsFileTag) *file.FileTagInfo {
	status := uint32(1)
	if row.Status != nil {
		status = uint32(*row.Status)
	}
	remark := ""
	if row.Remark != nil {
		remark = *row.Remark
	}
	return &file.FileTagInfo{
		Id:        ptr(uint64(row.ID)),
		CreatedAt: ptr(row.CreatedAt.UnixMilli()),
		UpdatedAt: ptr(row.UpdatedAt.UnixMilli()),
		Status:    ptr(status),
		Name:      ptr(row.Name),
		Remark:    ptr(remark),
	}
}

func toPBCloudFileTag(row *model.FmsCloudFileTag) *file.CloudFileTagInfo {
	status := uint32(1)
	if row.Status != nil {
		status = uint32(*row.Status)
	}
	remark := ""
	if row.Remark != nil {
		remark = *row.Remark
	}
	return &file.CloudFileTagInfo{
		Id:        ptr(uint64(row.ID)),
		CreatedAt: ptr(row.CreatedAt.UnixMilli()),
		UpdatedAt: ptr(row.UpdatedAt.UnixMilli()),
		Status:    ptr(status),
		Name:      ptr(row.Name),
		Remark:    ptr(remark),
	}
}

func toPBStorageProvider(row *model.FmsStorageProvider) *file.StorageProviderInfo {
	state := true
	if row.State != nil {
		state = *row.State
	}
	folder := ""
	if row.Folder != nil {
		folder = *row.Folder
	}
	cdnURL := ""
	if row.CdnURL != nil {
		cdnURL = *row.CdnURL
	}
	return &file.StorageProviderInfo{
		Id:        ptr(uint64(row.ID)),
		CreatedAt: ptr(row.CreatedAt.UnixMilli()),
		UpdatedAt: ptr(row.UpdatedAt.UnixMilli()),
		State:     ptr(state),
		Name:      ptr(row.Name),
		Bucket:    ptr(row.Bucket),
		SecretId:  ptr(row.SecretID),
		SecretKey: ptr(row.SecretKey),
		Endpoint:  ptr(row.Endpoint),
		Folder:    ptr(folder),
		Region:    ptr(row.Region),
		IsDefault: ptr(row.IsDefault),
		UseCdn:    ptr(row.UseCdn),
		CdnUrl:    ptr(cdnURL),
	}
}

func toPBCloudFile(row *model.FmsCloudFile, tagIDs []uint64) *file.CloudFileInfo {
	state := true
	if row.State != nil {
		state = *row.State
	}
	info := &file.CloudFileInfo{
		Id:          ptr(row.ID),
		CreatedAt:   ptr(row.CreatedAt.UnixMilli()),
		UpdatedAt:   ptr(row.UpdatedAt.UnixMilli()),
		State:       ptr(state),
		Name:        ptr(row.Name),
		Url:         ptr(row.URL),
		Size:        ptr(uint64(row.Size)),
		FileType:    ptr(uint32(row.FileType)),
		UserId:      ptr(row.UserID),
		TagIds:      tagIDs,
		RelativeSrc: ptr(""),
	}
	if row.CloudFileStorageProviders != nil {
		info.ProviderId = ptr(uint64(*row.CloudFileStorageProviders))
	}
	return info
}

func getFileTagIDs(ctx context.Context, q *query.Query, fileID string) []uint64 {
	rows, err := q.FileTagFile.WithContext(ctx).
		Select(q.FileTagFile.FileTagID).
		Where(q.FileTagFile.FileID.Eq(fileID)).Find()
	if err != nil {
		return nil
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, uint64(row.FileTagID))
	}
	return ids
}

func getFileTagIDsBatch(ctx context.Context, q *query.Query, fileIDs []string) map[string][]uint64 {
	result := make(map[string][]uint64, len(fileIDs))
	if len(fileIDs) == 0 {
		return result
	}
	rows, err := q.FileTagFile.WithContext(ctx).
		Select(q.FileTagFile.FileTagID, q.FileTagFile.FileID).
		Where(q.FileTagFile.FileID.In(fileIDs...)).Find()
	if err != nil {
		return result
	}
	for _, row := range rows {
		result[row.FileID] = append(result[row.FileID], uint64(row.FileTagID))
	}
	return result
}

func getCloudFileTagIDs(ctx context.Context, q *query.Query, fileID string) []uint64 {
	rows, err := q.CloudFileTagCloudFile.WithContext(ctx).
		Select(q.CloudFileTagCloudFile.CloudFileTagID).
		Where(q.CloudFileTagCloudFile.CloudFileID.Eq(fileID)).Find()
	if err != nil {
		return nil
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, uint64(row.CloudFileTagID))
	}
	return ids
}

func getCloudFileTagIDsBatch(ctx context.Context, q *query.Query, fileIDs []string) map[string][]uint64 {
	result := make(map[string][]uint64, len(fileIDs))
	if len(fileIDs) == 0 {
		return result
	}
	rows, err := q.CloudFileTagCloudFile.WithContext(ctx).
		Select(q.CloudFileTagCloudFile.CloudFileTagID, q.CloudFileTagCloudFile.CloudFileID).
		Where(q.CloudFileTagCloudFile.CloudFileID.In(fileIDs...)).Find()
	if err != nil {
		return result
	}
	for _, row := range rows {
		result[row.CloudFileID] = append(result[row.CloudFileID], uint64(row.CloudFileTagID))
	}
	return result
}

// fileIDsByTagIDs returns the ids of files carrying any of the given tags.
func fileIDsByTagIDs(ctx context.Context, q *query.Query, tagIDs []uint64) ([]string, error) {
	ids := make([]int64, 0, len(tagIDs))
	for _, id := range tagIDs {
		ids = append(ids, int64(id))
	}
	rows, err := q.FileTagFile.WithContext(ctx).
		Select(q.FileTagFile.FileID).
		Where(q.FileTagFile.FileTagID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	fileIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		fileIDs = append(fileIDs, row.FileID)
	}
	return fileIDs, nil
}

func cloudFileIDsByTagIDs(ctx context.Context, q *query.Query, tagIDs []uint64) ([]string, error) {
	ids := make([]int64, 0, len(tagIDs))
	for _, id := range tagIDs {
		ids = append(ids, int64(id))
	}
	rows, err := q.CloudFileTagCloudFile.WithContext(ctx).
		Select(q.CloudFileTagCloudFile.CloudFileID).
		Where(q.CloudFileTagCloudFile.CloudFileTagID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	fileIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		fileIDs = append(fileIDs, row.CloudFileID)
	}
	return fileIDs, nil
}
