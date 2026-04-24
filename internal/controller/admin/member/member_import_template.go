package member

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// memberImportTemplateBody UTF-8 BOM + 表头与示例行（密码需符合系统口令策略，示例为常见强密码格式）
const memberImportTemplateBody = "\xef\xbb\xbf用户名,密码,昵称,邮箱,手机,状态\r\ndemo_import,Aa1!demo88,示例昵称,,,0\r\n"

// ServeMemberImportTemplate 下载客户导入 CSV 模板（与 cmd 中注册路由一致）
func ServeMemberImportTemplate(r *ghttp.Request) {
	filename := "客户导入模板.csv"
	body := memberImportTemplateBody
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
