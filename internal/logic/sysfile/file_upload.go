package sysfile

import (
	"context"
	"io"
	"path/filepath"
	"strings"
	"time"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/utility/storage"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/google/uuid"
)

const maxSysFileUploadBytes = 100 << 20 // 100MB

// FileUpload 将文件写入当前活动存储并登记 sys_file_storage，返回新记录 id、存储 path、展示用文件名。
func (s *sSysfile) FileUpload(ctx context.Context, originalFilename string, size int64, contentType string, body io.ReadSeeker, isPrivate int, creator string) (id int64, objectPath string, displayName string, err error) {
	if size > maxSysFileUploadBytes {
		return 0, "", "", gerror.NewCode(consts.CodeInvalidParams)
	}
	if size <= 0 {
		return 0, "", "", gerror.NewCode(consts.CodeFileRequired)
	}
	stCfg, _ := storage.GetActiveConfig(ctx)
	adapter := storage.NewAdapter()
	bucket := strings.TrimSpace(stCfg.Bucket)
	if bucket == "" {
		bucket = "default"
	}
	base := filepath.Base(strings.ReplaceAll(originalFilename, "\\", "/"))
	if base == "" || base == "." {
		base = "unnamed"
	}
	ext := filepath.Ext(base)
	idPart := strings.ReplaceAll(uuid.New().String(), "-", "")
	objectKey := "uploads/" + time.Now().UTC().Format("2006/01") + "/" + idPart + ext
	if err := adapter.PutObject(ctx, bucket, objectKey, body, size, contentType); err != nil {
		return 0, "", "", gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if isPrivate != 0 && isPrivate != 1 {
		isPrivate = 0
	}
	r, err := dao.SysFileStorage.Ctx(ctx).Insert(sysdo.SysFileStorage{
		StorageType: stCfg.Type,
		Bucket:      bucket,
		Path:        objectKey,
		Filename:    base,
		MimeType:    contentType,
		Size:        size,
		Hash:        "",
		IsPrivate:   isPrivate,
		Creator:     creator,
		DeleteFlag:  consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		_ = adapter.Delete(ctx, bucket, objectKey)
		return 0, "", "", gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	newID, _ := r.LastInsertId()
	return newID, objectKey, base, nil
}

// FileOpenDownload 按文件记录打开存储对象，供 HTTP 下载流式写出。
func (s *sSysfile) FileOpenDownload(ctx context.Context, id int64) (filename string, mime string, size int64, body io.ReadCloser, err error) {
	var e sysentity.SysFileStorage
	err = dao.SysFileStorage.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil {
		return "", "", 0, nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if e.Id == 0 {
		return "", "", 0, nil, gerror.NewCode(consts.CodeFileNotFound)
	}
	adapter := storage.NewAdapter()
	rc, n, ct, err := adapter.GetObject(ctx, e.Bucket, e.Path)
	if err != nil {
		return "", "", 0, nil, gerror.WrapCode(consts.CodeFileNotFound, err, "")
	}
	mime = strings.TrimSpace(e.MimeType)
	if strings.TrimSpace(ct) != "" {
		mime = ct
	}
	if mime == "" {
		mime = "application/octet-stream"
	}
	size = e.Size
	if n > 0 {
		size = n
	}
	return e.Filename, mime, size, rc, nil
}
