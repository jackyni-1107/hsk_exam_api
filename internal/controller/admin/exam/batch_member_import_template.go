package exam

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

const batchMemberImportTemplateBody = "\xef\xbb\xbf用户名\r\ndemo_import\r\n"

// ServeBatchMemberImportTemplate 下载批次成员导入 CSV 模板。
func ServeBatchMemberImportTemplate(r *ghttp.Request) {
	filename := "批次成员导入模板.csv"
	body := batchMemberImportTemplateBody
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Length", strconv.Itoa(len(body)))
	safe := strings.Map(func(ch rune) rune {
		if ch >= 0x20 && ch < 0x7f && ch != '"' && ch != '\\' && ch != '/' {
			return ch
		}
		return '_'
	}, filename)
	if safe == "" {
		safe = "template.csv"
	}
	star := url.PathEscape(filename)
	r.Response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, safe, star))
	r.Response.WriteHeader(200)
	r.Response.Write([]byte(body))
}
