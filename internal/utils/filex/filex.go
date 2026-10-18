// Package filex contains file type/size/path helpers shared by the logic layer.
package filex

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/internal/config"
)

const (
	TypeOther uint32 = 1
	TypeImage uint32 = 2
	TypeVideo uint32 = 3
	TypeAudio uint32 = 4
)

// SplitMIME maps a MIME type (e.g. "image/png") to the reference file type
// string ("image"/"video"/"audio"/"other") and the persisted numeric type.
func SplitMIME(mime string) (string, uint32) {
	prefix := strings.SplitN(mime, "/", 2)[0]
	switch prefix {
	case "image":
		return "image", TypeImage
	case "video":
		return "video", TypeVideo
	case "audio":
		return "audio", TypeAudio
	default:
		return "other", TypeOther
	}
}

// CheckOverSize validates the file size against the per-type limit.
func CheckOverSize(fileType string, size int64, conf config.UploadConf) error {
	var limit int64
	switch fileType {
	case "image":
		limit = conf.MaxImageSize
	case "video":
		limit = conf.MaxVideoSize
	case "audio":
		limit = conf.MaxAudioSize
	default:
		limit = conf.MaxOtherSize
	}
	if limit > 0 && size > limit {
		return errorx.New(errorx.InvalidArgument, "file too large")
	}
	return nil
}

// SplitName splits a file name into base name and extension.
// ok is false when the name carries no dot (the reference rejects those).
func SplitName(name string) (base, ext string, ok bool) {
	dot := strings.LastIndex(name, ".")
	if dot == -1 {
		return "", "", false
	}
	return name[:dot], name[dot+1:], true
}

// LocalPath joins a store root with a DB relative path such as
// "/fms/image/2026-10-17/<uuid>.png".
func LocalPath(storePath, relPath string) string {
	return filepath.Join(storePath, filepath.FromSlash(strings.TrimPrefix(relPath, "/")))
}

// UUIDFromURL takes the last path segment of a URL and cuts it to the first
// 36 characters (the UUID), mirroring the reference delete_by_url behavior.
func UUIDFromURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", errorx.New(errorx.InvalidArgument, "wrong file path")
	}
	id := filepath.Base(u.Path)
	if len(id) < 36 {
		return "", errorx.New(errorx.InvalidArgument, "wrong file path")
	}
	return id[:36], nil
}
