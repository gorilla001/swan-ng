package mesos

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/andygrunwald/megos"
	"github.com/bbklab/swan-ng/mesos/protobuf/mesos"
	"github.com/samuel/go-zookeeper/zk"
)

// MesosState obtain current mesos stats
func (c *Client) MesosState() (*megos.State, error) {
	client, err := c.megosClient()
	if err != nil {
		return nil, err
	}
	return client.GetStateFromCluster()
}

// FrameworkState obtain current framework stats
func (c *Client) FrameworkState() (*megos.Framework, error) {
	stats, err := c.MesosState()
	if err != nil {
		return nil, err
	}

	fwName := c.framework.GetName()
	for _, fw := range stats.Frameworks {
		if fw.Name == fwName {
			nfw := fw
			return &nfw, nil
		}
	}

	return nil, fmt.Errorf("no such framework: %s", fwName)
}

// megosClient is just a helper mesos http client via vendor `andygrunwald/megos` which
// only `GET` on mesos http endpoints, we only use it to obtain cluster's states quickly.
func (c *Client) megosClient() (*megos.Client, error) {
	conn, connCh, err := zk.Connect(strings.Split(c.zkPath.Host, ","), 10*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// waiting for zookeeper to be connected.
	for event := range connCh {
		if event.State == zk.StateConnected {
			log.Info("connected to zookeeper succeed.")
			break
		}
	}

	var (
		masters    = make([]*url.URL, 0)
		masterInfo = new(mesos.MasterInfo)
	)

	children, _, err := conn.Children(c.zkPath.Path)
	if err != nil {
		return nil, fmt.Errorf("get children on %s error: %v", c.zkPath.Path, err)
	}

	for _, node := range children {
		if !strings.HasPrefix(node, "json.info") {
			break
		}

		path := c.zkPath.Path + "/" + node
		data, _, err := conn.Get(path)
		if err != nil {
			return nil, fmt.Errorf("get node on %s error: %v", path, err)
		}
		if err := json.Unmarshal(data, masterInfo); err != nil {
			return nil, err
		}

		var (
			ip   = *masterInfo.GetAddress().Ip
			port = *masterInfo.GetAddress().Port
		)
		masters = append(masters, &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", ip, port),
		})
	}

	return megos.NewClient(masters, nil), nil
}
