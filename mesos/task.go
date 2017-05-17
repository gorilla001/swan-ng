package mesos

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bbklab/swan-ng/mesos/protobuf/mesos"
	"github.com/bbklab/swan-ng/mesos/protobuf/sched"
	"github.com/bbklab/swan-ng/types"

	"github.com/gogo/protobuf/proto"
)

//
// utils about task op
//

type Task struct {
	ID          string  `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	IP          string  `json:"ip,omitempty"`
	Port        uint64  `json:"hostPorts,omitempty"`
	ContainerID string  `json:"containerId,omitempty"`
	Weight      float64 `json:"weight,omitempty"`
	State       string  `json:"state,omitempty"`
	ErrMsg      string  `json:"errMsg",omitempty"`
	CreatedAt   int64   `json:"createdAt,omitempty"`
	UpdatedAt   int64   `json:"updatedAt,omitempty"`
}

func (c *Client) LaunchTask(offer *mesos.Offer, task *mesos.TaskInfo) error {
	log.Println("launching task:", *task.Name)

	call := &sched.Call{
		FrameworkId: c.FrameworkId(),
		Type:        sched.Call_ACCEPT.Enum(),
		Accept: &sched.Call_Accept{
			OfferIds: []*mesos.OfferID{
				offer.GetId(),
			},
			Operations: []*mesos.Offer_Operation{
				&mesos.Offer_Operation{
					// TODO replace with LAUNCH_GROUP
					Type: mesos.Offer_Operation_LAUNCH.Enum(),
					Launch: &mesos.Offer_Operation_Launch{
						TaskInfos: []*mesos.TaskInfo{task},
					},
				},
			},
			Filters: &mesos.Filters{RefuseSeconds: proto.Float64(1)},
		},
	}

	// send call
	err := c.Send(call)
	if err != nil {
		return err
	}

	// subcribe waitting for task's update events here until task finished or met error.
	for {
		update := c.SubscribeTaskUpdate(task.TaskId.GetValue())

		// for debug
		//json.NewEncoder(os.Stdout).Encode(update)

		status := update.GetUpdate().GetStatus()
		err = c.AckUpdateEvent(status)
		if err != nil {
			break
		}

		if c.IsTaskDone(status) {
			err = c.DetectError(status) // check if we met an error.
			break
		}
	}

	return err
}

func (c *Client) UpdateTask(task *Task) error {
	return nil
}

func (c *Client) KillTask(taskID, agentID string) error {
	log.Println("stopping task:", taskID)

	call := &sched.Call{
		FrameworkId: c.FrameworkId(),
		Type:        sched.Call_KILL.Enum(),
		Kill: &sched.Call_Kill{
			TaskId: &mesos.TaskID{
				Value: proto.String(taskID),
			},
			AgentId: &mesos.AgentID{
				Value: proto.String(agentID),
			},
		},
	}

	// send call
	err := c.Send(call)
	if err != nil {
		return err
	}

	// subcribe waitting for task's update events here until task finished or met error.
	for {
		update := c.SubscribeTaskUpdate(taskID)

		// for debug
		//json.NewEncoder(os.Stdout).Encode(update)

		status := update.GetUpdate().GetStatus()
		err = c.AckUpdateEvent(status)
		if err != nil {
			break
		}

		if c.IsTaskDone(status) {
			err = c.DetectError(status) // check if we met an error.
			break
		}
	}

	return nil
}

type taskOption struct {
	id      string
	name    string
	version *types.AppVersion
}

func (to *taskOption) newTask(offer *mesos.Offer) *mesos.TaskInfo {
	task := &mesos.TaskInfo{
		Name:        to.taskName(),
		TaskId:      to.taskID(),
		AgentId:     offer.GetAgentId(),
		Resources:   to.resourceInfo(),
		Command:     to.commandInfo(),
		Container:   to.containerInfo(),
		HealthCheck: to.version.HealthCheck,
		KillPolicy:  to.version.KillPolicy,
		Labels:      to.labels(),
	}

	return task
}

func (to *taskOption) taskID() *mesos.TaskID {
	return &mesos.TaskID{Value: proto.String(to.id)}
}

func (to *taskOption) taskName() *string {
	return proto.String(to.Name)
}

func (to *taskOption) commandInfo() *mesos.CommandInfo {
	if cmd := to.service.Service.Command; len(cmd) > 0 {
		return &mesos.CommandInfo{
			Shell:     proto.Bool(false),
			Value:     proto.String(cmd[0]),
			Arguments: cmd[1:],
		}
	}
	return &mesos.CommandInfo{Shell: proto.Bool(false)}
}

func (to *taskOption) containerInfo() *mesos.ContainerInfo {
	var (
		image      = to.version.Container.Docker.Image
		privileged = to.version.Container.Docker.Privileged
		force      = to.version.Container.Docker.ForcePullImage
	)

	return &mesos.ContainerInfo{
		Type:    mesos.ContainerInfo_DOCKER.Enum(),
		Volumes: to.volumes(),
		Docker: &mesos.ContainerInfo_DockerInfo{
			Image:          image,
			Privileged:     privileged,
			Network:        to.network(),
			PortMappings:   to.portMappings(),
			Parameters:     to.paramters(),
			ForcePullImage: force,
		},
	}
}

func (to *taskOption) resourceInfo(offer *mesos.Offer) *mesos.Resource {
	var (
		v   = to.version
		pms = v.Container.Docker.PortMappings
		rs  = make([]*mesos.Resource, 0, 0)
	)

	if v.Cpus > 0 {
		rs = append(rs, &mesos.Resource{
			Name: proto.String("cpus"),
			Type: mesos.Value_SCALAR.Enum(),
			Scalar: &mesos.Value_Scalar{
				Value: proto.Float64(ro.CPU),
			},
		})
	}

	if v.Mem > 0 {
		rs = append(rs, &mesos.Resource{
			Name: proto.String("mem"),
			Type: mesos.Value_SCALAR.Enum(),
			Scalar: &mesos.Value_Scalar{
				Value: proto.Float64(v.Mem),
			},
		})
	}

	if ro.Disk > 0 {
		rs = append(rs, &mesos.Resource{
			Name: proto.String("disk"),
			Type: mesos.Value_SCALAR.Enum(),
			Scalar: &mesos.Value_Scalar{
				Value: proto.Float64(v.Disk),
			},
		})
	}

	for i, m := range pms {
		ports := GetPorts(offer)

		rs = append(rs, &mesos.Resource{
			Name: proto.String("ports"),
			Type: mesos.Value_RANGES.Enum(),
			Ranges: &mesos.Value_Ranges{
				Range: []*mesos.Value_Range{
					{
						Begin: proto.Uint64(uint64(ports[i])),
						End:   proto.Uint64(uint64(ports[i])),
					},
				},
			},
		})
	}

	return rs
}

func (to *taskOption) labels() *mesos.Labels {
	return to.version.Labels
}

func (to *taskOption) containerInfo(offer *mesos.Offer) *mesos.ContainerInfo {
	var (
		image      = to.version.Container.Docker.Image
		privileged = to.version.Container.Docker.Privileged
		parameters = to.version.Container.Docker.Parameters
		volumes    = to.version.Container.Volumes
		network    = to.version.Container.Docker.Network
	)

	info := &mesos.ContainerInfo{
		Type: mesos.ContainerInfo_DOCKER.Enum(),
		Docker: &mesos.ContainerInfo_DockerInfo{
			Image:      proto.String(image),
			Privileged: proto.Bool(privileged),
		},
	}
}
func (to *taskOption) parameters() []*mesos.Parameter {
	var (
		parameters = to.version.Container.Docker.Parameters
		mps        = make([]*mesos.Parameter, 0, 0)
	)

	for _, para := range parameters {
		mps = append(mps, &mesos.Parameter{
			Key:   proto.String(para.Key),
			Value: proto.String(para.Value),
		})
	}

	return mps
}

func (to *taskOption) volumes() []*mesos.Volumes {
	var (
		vlms = to.version.Container.Volumes
		mvs  = make([]*mesos.Volume, 0, 0)
	)

	for _, vlm := range vlms {
		mode := mesos.Volume_RO
		if vlm.Mode == "RW" {
			mode = mesos.Volume_RW
		}

		mvs = append(mvs, &mesos.Volume{
			ContainerPath: proto.String(vlm.ContainerPath),
			HostPath:      proto.String(vlm.HostPath),
			Mode:          &mode,
		})
	}

	return mvs
}

func (to *taskOption) network() *mesos.ContainerInfo_DockerInfo_Network {
	var netMode mesos.ContainerInfo_DockerInfo_Network

	switch net := to.version.Container.Docker.Network; net {
	case "none":
		netMode = mesos.ContainerInfo_DockerInfo_NONE.Enum()
	case "host":
		netMode = mesos.ContainerInfo_DockerInfo_HOST.Enum()
	case "bridge":
		netMode = mesos.ContainerInfo_DockerInfo_BRIDGE.Enum()
	case "user":
		netMode = mesos.ContainerInfo_DockerInfo_USER.Enum()
	default:
		netMode = mesos.ContainerInfo_DockerInfo_NONE.Enum()
	}

	return netMode
}

func (to *taskOption) portMappings() []*mesos.ContainerInfo_DockerInfo_PortMapping {
	var (
		pms  = to.version.Container.Docker.PortMappings
		dpms = make([]*mesos.ContainerInfo_DockerInfo_PortMapping, 0, 0)
	)

	ports := to.GetPorts(offer)
	for i, m := range pms {
		dpms = append(dpms,
			&mesos.ContainerInfo_DockerInfo_PortMapping{
				HostPort:      proto.Uint32(uint32(ports[i])),
				ContainerPort: proto.Uint32(uint32(m.ContainerPort)),
				Protocol:      proto.String(m.Protocol),
			})
	}

	return dpms
}

func (to *taskOption) GetPorts(offer *mesos.Offer) (ports []uint64) {
	for _, resource := range offer.Resources {
		if resource.GetName() == "ports" {
			for _, rang := range resource.GetRanges().GetRange() {
				for i := rang.GetBegin(); i <= rang.GetEnd(); i++ {
					ports = append(ports, i)
				}
			}
		}
	}
	return ports
}

func (c *Client) AckUpdateEvent(taskStatus *mesos.TaskStatus) error {
	if taskStatus.GetUuid() != nil {
		call := &sched.Call{
			FrameworkId: c.FrameworkId(),
			Type:        sched.Call_ACKNOWLEDGE.Enum(),
			Acknowledge: &sched.Call_Acknowledge{
				AgentId: taskStatus.GetAgentId(),
				TaskId:  taskStatus.GetTaskId(),
				Uuid:    taskStatus.GetUuid(),
			},
		}
		return c.Send(call)
	}
	return nil
}
