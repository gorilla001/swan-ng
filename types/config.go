package types

import (
	"errors"
	"fmt"
	"net/url"
)

// MgrConfig represents manager configs
type MgrConfig struct {
	ListenAddr  string   `json:"listen_addr"`
	MesosZKPath *url.URL `json:"mesos_zk_path"` // mesos zk addr
	ZKPath      *url.URL `json:"zk_path"`       // swan zk store addr, if null, use memory store
}

// Valid verify the manager configs
func (c *MgrConfig) Valid() error {
	if c.ListenAddr == "" {
		return errors.New("listen param required")
	}

	p := c.MesosZKPath
	if p == nil {
		return errors.New("mesos zk_path param required")
	}
	if err := validZKURL(p); err != nil {
		return fmt.Errorf("mesos zk path invalid: %v", err)
	}

	if p := c.ZKPath; p != nil {
		if err := validZKURL(p); err != nil {
			return fmt.Errorf("swan zk path invalid: %v", err)
		}
	}

	return nil
}

func validZKURL(url *url.URL) error {
	if url.Host == "" {
		return fmt.Errorf("url Host required")
	}
	if url.Scheme != "zk" {
		return fmt.Errorf("url Scheme should be zk://")
	}
	if url.Path == "" {
		return fmt.Errorf("url Path required")
	}
	return nil
}
