// Package types ...
// NOTE: for compatibility
// most of these definations are copied from original swan store/structs.go
package types

// App ...
type App struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	ClusterID string `json:"clusterId,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
	State     string `json:"state,omitempty"`
	// StateMachine *StateMachine `json:"stateMachine,omitempty"`

	// app settings
	Version         *AppVersion `json:"version,omitempty"`
	ProposedVersion *AppVersion `json:"proposedVersion,omitempty"`
}

// AppWrapper is only for display, it wraps `App` with more useful fields.
// TODO sigh, for compatibility, should keep same as original swan types/app.go
type AppWrapper struct {
	*App
	//UpdatedInstances int      `json:"updatedInstances"`
	//RunningInstances int      `json:"runningInstances"`
	//Tasks            []*Task  `json:"tasks,omitempty"`
	//Versions         []string `json:"versions,omitempty"`
}

// AppVersion ...
type AppVersion struct {
	//ID           string            `json:"id,omitempty"`
	//AppID        string            `json:"appID,omitempty"`
	AppName      string            `json:"appName,omitempty"`
	AppVersion   string            `json:"appVersion,omitempty"`
	Command      string            `json:"command,omitempty"`
	Cpus         float64           `json:"cpus,omitempty"`
	Mem          float64           `json:"mem,omitempty"`
	Disk         float64           `json:"disk,omitempty"`
	Instances    int32             `json:"instances,omitempty"`
	RunAs        string            `json:"runAs,omitempty"`
	Container    *Container        `json:"container,omitempty"`
	Labels       map[string]string `json:"labels"`
	HealthCheck  *HealthCheck      `json:"healthCheck,omitempty"`
	Env          map[string]string `json:"env"`
	KillPolicy   *KillPolicy       `json:"killPolicy,omitempty"`
	UpdatePolicy *UpdatePolicy     `json:"updatePolicy,omitempty"`
	Gateway      *Gateway          `json:"gateway,omitempty"`
	Constraints  string            `json:"constraints,omitempty"`
	Uris         []string          `json:"uris,omitempty"`
	IP           []string          `json:"ip,omitempty"`
	Mode         string            `json:"mode,omitempty"`
	Priority     int32             `json:"priority,omitempty"`
}

// Container ...
type Container struct {
	Type    string    `json:"type,omitempty"`
	Docker  *Docker   `json:"docker,omitempty"`
	Volumes []*Volume `json:"volumes,omitempty"`
}

// Docker ...
type Docker struct {
	ForcePullImage bool           `json:"forcePullImage,omitempty"`
	Image          string         `json:"image,omitempty"`
	Network        string         `json:"network,omitempty"`
	Parameters     []*Parameter   `json:"parameters,omitempty"`
	PortMappings   []*PortMapping `json:"portMappings,omitempty"`
	Privileged     bool           `json:"privileged,omitempty"`
}

// Parameter ...
type Parameter struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// PortMapping ...
type PortMapping struct {
	ContainerPort int32  `json:"containerPort,omitempty"`
	HostPort      int32  `json:"hostPort,omitempty"`
	Name          string `json:"name,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

// Volume ...
type Volume struct {
	ContainerPath string `json:"containerPath,omitempty"`
	HostPath      string `json:"hostPath,omitempty"`
	Mode          string `json:"mode,omitempty"`
}

// KillPolicy ...
type KillPolicy struct {
	Duration int64 `json:"duration,omitempty"`
}

// UpdatePolicy ...
type UpdatePolicy struct {
	UpdateDelay  int32  `json:"updateDelay,omitempty"`
	MaxRetries   int32  `json:"maxRetries,omitempty"`
	MaxFailovers int32  `json:"maxFailovers,omitempty"`
	Action       string `json:"action,omitempty"`
}

// Gateway ...
type Gateway struct {
	Enabled bool    `json:"enabled,omitempty"`
	Weight  float64 `json:"weight,omitempty"`
}

// HealthCheck ...
type HealthCheck struct {
	ID                  string  `json:"id,omitempty"`
	Address             string  `json:"address,omitempty"`
	Protocol            string  `json:"protocol,omitempty"`
	Port                int32   `json:"port,omitempty"`
	PortIndex           int32   `json:"portIndex,omitempty"`
	PortName            string  `json:"portName,omitempty"`
	Value               string  `json:"value,omitempty"`
	Path                string  `json:"path,omitempty"`
	ConsecutiveFailures uint32  `json:"consecutiveFailures,omitempty"`
	GracePeriodSeconds  float64 `json:"gracePeriodSeconds,omitempty"`
	IntervalSeconds     float64 `json:"intervalSeconds,omitempty"`
	TimeoutSeconds      float64 `json:"timeoutSeconds,omitempty"`
	DelaySeconds        float64 `json:"delaySeconds,omitempty"`
}

// Task ...
type Task struct {
	ID            string   `json:"id,omitempty"`
	AppID         string   `json:"appId,omitempty"`
	VersionID     string   `json:"versionId,omitempty"`
	State         string   `json:"state,omitempty"`
	Stdout        string   `json:"stdout,omitempty"`
	Stderr        string   `json:"stderr,omitempty"`
	HostPorts     []uint64 `json:"hostPorts,omitempty"`
	OfferID       string   `json:"offerId,omitempty"`
	AgentID       string   `json:"agentId,omitempty"`
	IP            string   `json:"ip,omitempty"`
	AgentHostName string   `json:"agentHostName,omitempty"`
	Reason        string   `json:"reason,omitempty"`
	Message       string   `json:"message,omitempty"`
	CreatedAt     int64    `json:"createdAt,omitempty"`
	ArchivedAt    int64    `json:"archivedAt,omitempty"`
	ContainerID   string   `json:"containerId,omitempty"`
	ContainerName string   `json:"containerName,omitempty"`
	Weight        float64  `json:"weight,omitempty"`
	//SlotID        string   `json:"slotId,omitempty"`
}
