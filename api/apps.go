package api

import (
	"encoding/json"
	"fmt"

	"github.com/bbklab/swan-ng/api/mux"
	"github.com/bbklab/swan-ng/store"
	"github.com/bbklab/swan-ng/types"
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

func createApp(ctx *mux.Context) {
	dec := json.NewDecoder(ctx.Req.Body)
	ver := new(types.AppVersion)
	if err := dec.Decode(&ver); err != nil {
		fmt.Println(err)
	}

	if err := mesosCli.LaunchApp(ver); err != nil {
		ctx.Error(500, err)
	}

	ctx.JSON(201, "succeed")
}
