package mesos

import (
	"os"

	"github.com/golang/protobuf/proto"

	"github.com/bbklab/swan-ng/mesos/protobuf/mesos"
)

func defaultFramework() *mesos.FrameworkInfo {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "UNKNOWN"
	}

	return &mesos.FrameworkInfo{
		// ID:              proto.String(""), // reset later
		User:            proto.String("root"),
		Name:            proto.String("swan"),
		Principal:       proto.String("swan"),
		FailoverTimeout: proto.Float64(60 * 60 * 24 * 7),
		Checkpoint:      proto.Bool(false),
		Hostname:        proto.String(hostName),
		Capabilities: []*mesos.FrameworkInfo_Capability{
			{Type: mesos.FrameworkInfo_Capability_PARTITION_AWARE.Enum()},
			{Type: mesos.FrameworkInfo_Capability_TASK_KILLING_STATE.Enum()},
		},
	}
}
