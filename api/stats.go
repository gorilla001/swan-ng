package api

import (
	"strings"

	"github.com/bbklab/swan-ng/api/mux"
	"github.com/bbklab/swan-ng/store"
	"github.com/bbklab/swan-ng/types"
)

// GET /stats
func stats(ctx *mux.Context) {
	ms, err := mesosCli.MesosState()
	if err != nil {
		ctx.Error(500, err)
		return
	}

	ret := &types.Stats{
		ClusterID: mesosCli.Cluster(),
		Created:   ms.StartTime,
		Master:    strings.Split(ms.Leader, "@")[1],
		AppStats:  make(map[string]int),
	}

	apps, err := store.DB().ListApps()
	if err != nil {
		ctx.Error(500, err)
	}

	for _, app := range apps {
		ret.AppCount++
		tasks, _ := store.DB().ListTasks(app.ID)
		ret.TaskCount += len(tasks)
		ret.AppStats[app.Version.RunAs]++
	}

	ss := make([]string, 0, len(ms.Slaves))
	for _, slave := range ms.Slaves {
		// TODO verify ...
		ret.TotalCPU += slave.Resources.CPUs
		ret.TotalMem += slave.Resources.Mem
		ret.TotalDisk += slave.Resources.Disk
		ret.CPUTotalUsed += slave.UsedResources.CPUs
		ret.MemTotalUsed += slave.UsedResources.CPUs
		ret.DiskTotalUsed += slave.UsedResources.CPUs
		ret.CPUTotalOffered += slave.OfferedResources.CPUs
		ret.MemTotalOffered += slave.OfferedResources.Mem
		ret.DiskTotalOffered += slave.OfferedResources.Disk

		if attrs := slave.Attributes; len(attrs) != 0 {
			ret.Attributes = append(ret.Attributes, attrs)
		}

		ss = append(ss, strings.Split(slave.PID, "@")[1])
	}
	ret.Slaves = strings.Join(ss, ",")

	ctx.JSON(200, ret)
}
