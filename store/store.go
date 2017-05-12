package store

import (
	"fmt"
	"net/url"

	"github.com/bbklab/swan-ng/store/memory"
	"github.com/bbklab/swan-ng/store/zk"
	"github.com/bbklab/swan-ng/types"
)

var (
	db Store
)

// Setup initizliz specified type of db store
func Setup(typ string, url *url.URL) (err error) {
	switch typ {
	case "memory":
		db = memory.New()
	case "zk":
		db, err = zk.New(url)
	default:
		err = fmt.Errorf("unsupported db store type: %v", typ)
	}

	return
}

// DB return the global store.Store implement
func DB() Store {
	return db
}

// Store interface defination
type Store interface {
	// framework
	GetFrameworkID() string
	UpdateFrameworkID(id string) error

	// app CRUD
	CreateApp(app *types.App) error
	UpdateApp(app *types.App) error
	GetApp(id string) (*types.App, error)
	ListApps() ([]*types.App, error)
	DeleteApp(id string) error
	// app's setting version
	CreateVersion(aid string, ver *types.Version) error
	GetVersion(aid, vid string) (*types.Version, error)
	ListVersions(aid string) ([]*types.Version, error)
	// app's tasks
	UpdateTask(aid string, t *types.Task) error              // update app's specified task
	ListTasks(aid string) ([]*types.Task, error)             // app's task list
	GetTaskHistories(aid, tid string) ([]*types.Task, error) // app's specified task's histories

	// compose instance CRUD
	CreateInstance(ins *types.Instance) error
	UpdateInstance(ins *types.Instance) error // status, errmsg, updateAt
	GetInstance(id string) (*types.Instance, error)
	ListInstances() ([]*types.Instance, error)
	DeleteInstance(id string) error
}
