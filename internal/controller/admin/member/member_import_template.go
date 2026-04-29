package member

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// memberImportTemplateBody UTF-8 BOM + 表头与示例行（无状态列；密码可省略，省略时由邮箱第 1、3、5 位 + @hskmock 生成）
const memberImportTemplateBody = "\xef\xbb\xbf昵称,邮箱,手机\r\n示例昵称,demo@example.com,\r\n"

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
