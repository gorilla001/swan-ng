package api

import (
	"github.com/bbklab/swan-ng/api/mux"
	"github.com/bbklab/swan-ng/store"
)

// GET /apps
func listApps(ctx *mux.Context) {
	apps, err := store.DB().ListApps()
	if err != nil {
		ctx.Error(500, err)
		return
	}

	ctx.JSON(200, apps)
}

// GET /apps/:id
func getApp(ctx *mux.Context) {
	id := ctx.Ps["id"]

	app, err := store.DB().GetApp(id)
	if err != nil {
		ctx.Error(500, err)
		return
	}

	// TODO wrap app within types.AppWrapper
	ctx.JSON(200, app)
}
