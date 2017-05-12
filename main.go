package main

import (
	"net/url"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/bbklab/swan-ng/api"
	"github.com/bbklab/swan-ng/types"
	"github.com/bbklab/swan-ng/version"
)

var (
	globalFlags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "DEBUG",
		},
	}

	managerFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "listen-addr",
			Usage:  "http listener address",
			EnvVar: "SWAN_LISTEN_ADDR",
			Value:  "0.0.0.0:9999",
		},
		cli.StringFlag{
			Name:   "mesos",
			Usage:  "mesos zookeeper path. eg. zk://host1:port1,host2:port2,.../mesos",
			EnvVar: "SWAN_MESOS_URL",
		},
		cli.StringFlag{
			Name:   "zk",
			Usage:  "swan zookeeper path. eg. zk://host1:port1,host2:port2,.../swan",
			EnvVar: "SWAN_ZK_URL",
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "swan"
	app.Version = version.GetVersion()
	if gitCommit := version.GetGitCommit(); gitCommit != "" {
		app.Version += "-" + gitCommit
	}

	app.Flags = globalFlags

	app.Before = func(c *cli.Context) error {
		debug := c.Bool("debug")

		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		log.SetLevel(log.InfoLevel)
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		log.SetOutput(os.Stdout)
		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "manager",
			ShortName: "m",
			Usage:     "run as manager",
			Flags:     managerFlags,
			Action: func(c *cli.Context) error {
				cfg, err := newMgrCfg(c)
				if err != nil {
					return err
				}
				return api.Serve(cfg)
			},
		},
		cli.Command{
			Name:      "version",
			ShortName: "v",
			Usage:     "print version",
			Action: func(c *cli.Context) error {
				return version.Version().FormatTo(os.Stdout)
			},
		},
	}

	app.RunAndExitOnError()
}

func newMgrCfg(c *cli.Context) (*types.MgrConfig, error) {
	var (
		listen  = c.String("listen-addr")
		mesosZK = c.String("mesos")
		zk      = c.String("zk")
		err     error
	)

	cfg := &types.MgrConfig{
		ListenAddr: listen,
	}

	if cfg.MesosZKPath, err = url.Parse(mesosZK); err != nil {
		return nil, err
	}

	if zk != "" { // allow null, if null, use memory store
		if cfg.ZKPath, err = url.Parse(zk); err != nil {
			return nil, err
		}
	}

	return cfg, cfg.Valid()
}
