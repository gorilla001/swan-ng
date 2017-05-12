package api

import (
	"github.com/bbklab/swan-ng/api/mux"
	"github.com/bbklab/swan-ng/version"
)

// GET /version
func showVersion(ctx *mux.Context) {
	ctx.JSON(200, version.Version())
}
