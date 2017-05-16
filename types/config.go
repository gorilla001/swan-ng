package types

import (
	"fmt"
	"net/url"
)

// MgrConfig represents manager configs
type MgrConfig struct {
	Listen   string   `json:"listen"`
	MesosURL *url.URL `json:"mesos"` // mesos zk addr
	ZKURL    *url.URL `json:"zk"`    // swan zk store addr, if null, use memory store
}

// Valid verify the manager configs
func (c *MgrConfig) Valid() error {
	if c.Listen == "" {
		return fmt.Errorf("listen param required")
	}

	if err := validZKURL(c.MesosURL); err != nil {
		return fmt.Errorf("mesos zk url invalid: %v", err)
	}

	if p := c.ZKURL; p != nil {
		if err := validZKURL(p); err != nil {
			return fmt.Errorf("swan zk url invalid: %v", err)
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
