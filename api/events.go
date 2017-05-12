package api

import (
	"net/http"

	"github.com/bbklab/swan-ng/api/mux"
)

// GET /events
func events(ctx *mux.Context) {
	ctx.Res.Header().Set("Content-Type", "text/event-stream")
	ctx.Res.Header().Set("Cache-Control", "no-cache")
	ctx.Res.Write(nil)

	if f, ok := ctx.Res.(http.Flusher); ok {
		f.Flush()
	}

	if eventMgr.full() {
		ctx.Error(500, "too many event clients")
		return
	}

	remote := ctx.Req.RemoteAddr
	eventMgr.subscribe(remote, ctx.Res)
	eventMgr.wait(remote)
}
