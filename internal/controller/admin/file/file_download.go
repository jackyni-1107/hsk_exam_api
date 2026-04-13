package file

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"exam/internal/consts"
	sysfilesvc "exam/internal/service/sysfile"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ServeDownload 流式下载文件（不经 JSON 包装中间件），需与 cmd 中注册的路由一致。
func ServeDownload(r *ghttp.Request) {
	ctx := r.GetCtx()
	id := r.Get("id").Int64()
	if id <= 0 {
		writeDownloadJSONErr(r, 400, "invalid id")
		return
	}
	filename, mime, size, rc, err := sysfilesvc.Sysfile().FileOpenDownload(ctx, id)
	if err != nil {
		st := 500
		if c := gerror.Code(err); c != nil && c.Code() == consts.CodeFileNotFound.Code() {
			st = 404
		}
		msg := "error"
		if c := gerror.Code(err); c != nil {
			msg = c.Message()
		}
		writeDownloadJSONErr(r, st, msg)
		return
	}
	defer rc.Close()

	r.Response.Header().Set("Content-Type", mime)
	if size > 0 {
		r.Response.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	}
	safe := strings.Map(func(r rune) rune {
		if r >= 0x20 && r < 0x7f && r != '"' && r != '\\' && r != '/' {
			return r
		}
		return '_'
	}, filename)
	if safe == "" {
		safe = "download"
	}
	star := url.PathEscape(filename)
	r.Response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, safe, star))
	r.Response.WriteHeader(200)
	_, _ = io.Copy(r.Response.Writer, rc)
}

func writeDownloadJSONErr(r *ghttp.Request, code int, msg string) {
	r.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.Response.WriteHeader(code)
	b, _ := json.Marshal(map[string]interface{}{"code": code, "message": msg})
	r.Response.Write(b)
}
