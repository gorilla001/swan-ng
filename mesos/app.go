package mesos

import (
	"fmt"
	"time"

	"github.com/bbklab/swan-ng/types"

	uuid "github.com/satori/go.uuid"
)

const (
	StatusCreating = "creating"
)

type App struct {
	ID              string            `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	ClusterID       string            `json:"clusterId,omitempty"`
	State           string            `json:"state,omitempty"`
	Version         *types.AppVersion `json:"version,omitempty"`
	ProposedVersion *types.AppVersion `json:"proposedVersion,omitempty"`
	CreatedAt       time.Time         `json:"createdAt,omitempty"`
	UpdatedAt       time.Time         `json:"updatedAt,omitempty"`
}

func (c *Client) LaunchApp(version *types.AppVersion) error {
	id := fmt.Sprintf("%s.%s.%s", version.AppName, version.RunAs, c.Cluster())
	app := &App{
		ID:              id,
		Name:            version.AppName,
		ClusterID:       c.Cluster(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		State:           StatusCreating,
		Version:         version,
		ProposedVersion: nil,
	}

	for i := 0; i < int(version.Instances); i++ {
		res := Resource{
			mem:   version.Mem,
			cpus:  version.Cpus,
			disk:  version.Disk,
			ports: len(version.Container.Docker.PortMappings),
		}

		offer := c.SubscribeOffer(res)

		to := &taskOption{
			fmt.Sprintf("%s.%s", i, app.ID),
			fmt.Sprintf("%s.%s", i, app.ID),
			version,
		}

		task := to.newTask(offer)

		if err := c.LaunchTask(offer, task); err != nil {
			return err
		}

	}

	return nil
}
