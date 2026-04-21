package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			InitAll(ctx)

			s := g.Server()
			registerHTTPRoutes(s)
			startBackgroundRuntimes(ctx)
			s.Run()
			return nil
		},
	}
)
