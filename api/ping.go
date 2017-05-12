package api

import "github.com/bbklab/swan-ng/api/mux"

// GET /ping
func ping(ctx *mux.Context) {
	ctx.Text(200, "pong")
}
