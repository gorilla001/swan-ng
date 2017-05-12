package types

// Stats ...
// NOTE copyid from original swan src/types
type Stats struct {
	ClusterID  string                   `json:"clusterID"`
	AppCount   int                      `json:"appCount"`
	TaskCount  int                      `json:"taskCount"`
	Created    float64                  `json:"created"`
	Master     string                   `json:"master"`
	Slaves     string                   `json:"slaves"`
	Attributes []map[string]interface{} `json:"attributes"`
	AppStats   map[string]int           `json:"appStats"` // runas -> nb

	// resource usages
	TotalCPU         float64 `json:"totalCpu"`
	TotalMem         float64 `json:"totalMem"`
	TotalDisk        float64 `json:"totalDisk"`
	CPUTotalOffered  float64 `json:"cpuTotalOffered"`
	MemTotalOffered  float64 `json:"memTotalOffered"`
	DiskTotalOffered float64 `json:"diskTotalOffered"`
	CPUTotalUsed     float64 `json:"cpuTotalUsed"`
	MemTotalUsed     float64 `json:"memTotalUsed"`
	DiskTotalUsed    float64 `json:"diskTotalUsed"`
}
