package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bbklab/swan-ng/api/mux"
	"github.com/bbklab/swan-ng/mesos"
	"github.com/bbklab/swan-ng/store"
	"github.com/bbklab/swan-ng/types"
)

const (
	// APIPREFIX is exported
	APIPREFIX = "/api"
)

// package scope instances
var (
	cfg      = new(types.MgrConfig) // manager configs
	mesosCli = new(mesos.Client)    // mesos scheduler client
)

// Serve initilize & startup http api services
func Serve(mc *types.MgrConfig) (err error) {
	// setup manager configs
	if mc == nil {
		return errors.New("manager config requried")
	}
	*cfg = *mc

	// setup mesos client & startup mesos events subscriber
	mesosCli, err = mesos.NewClient(mc.MesosURL)
	if err != nil {
		return fmt.Errorf("initialize mesos client error: [%v]", err)
	}
	if err := mesosCli.Subscribe(); err != nil {
		return fmt.Errorf("startup mesos events subscriber error: [%v]", err)
	}

	// setup swan db store
	if url := cfg.ZKURL; url == nil {
		err = store.Setup("memory", nil)
	} else {
		err = store.Setup("zk", url)
	}
	if err != nil {
		return fmt.Errorf("initialize db store error: [%v]", err)
	}

	// setup http routes & serving
	mux := mux.New()
	setupRouters(mux)
	server := http.Server{
		Addr:    cfg.Listen,
		Handler: mux,
	}

	return server.ListenAndServe()
}

func setupRouters(m *mux.Mux) {
	m.Get("/", listAPIs)

	m.Get("/ping", ping)
	m.Get("/events", events)
	m.Get("/stats", stats)
	m.Get("/version", showVersion)

	// apps
	m.Get("/apps", listApps)
	m.Get("/apps/:id", getApp)
	//m.Delete("/apps/:id", delApp)
	//m.Patch("/apps/:id/scale", scaleApp)
}

// GET /
func listAPIs(ctx *mux.Context) {
	rs := ctx.Mux.AllRoutes()
	ss := make([]string, 0, len(rs))
	for _, r := range rs {
		ss = append(ss, fmt.Sprintf("%s %s", r.Method, r.Pattern))
	}
	ctx.JSON(200, ss)
}
