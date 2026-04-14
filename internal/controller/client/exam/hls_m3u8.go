package exam

import (
	papersvc "exam/internal/service/paper"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ServeHlsM3U8 返回原始 m3u8（不经 JSON 包装）；ticket 路径段可含 .m3u8 后缀。
func ServeHlsM3U8(r *ghttp.Request) {
	ticket := r.Get("ticket").String()
	body, err := papersvc.Paper().BuildHlsM3U8Playlist(r.GetCtx(), ticket)
	if err != nil {
		// 打印錯誤日志
		g.Log().Errorf(r.GetCtx(), "ServeHlsM3U8 error: %v", err)
		r.Response.WriteHeader(404)
		return
	}
	r.Response.Header().Set("Content-Type", "application/vnd.apple.mpegurl; charset=utf-8")
	r.Response.Header().Set("Cache-Control", "no-store")
	r.Response.Write(body)
}
