package health

import (
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Liveness 进程存活（不探测依赖）
func Liveness(r *ghttp.Request) {
	r.Response.WriteStatus(http.StatusOK, `{"status":"ok"}`)
}

// Readiness 探测 MySQL 与 Redis
func Readiness(r *ghttp.Request) {
	ctx := r.GetCtx()
	if err := g.DB().PingMaster(); err != nil {
		r.Response.WriteStatus(http.StatusServiceUnavailable, `{"status":"unready","db":"down"}`)
		return
	}
	if _, err := g.Redis().Do(ctx, "PING"); err != nil {
		r.Response.WriteStatus(http.StatusServiceUnavailable, `{"status":"unready","redis":"down"}`)
		return
	}
	r.Response.WriteStatus(http.StatusOK, `{"status":"ready"}`)
}
