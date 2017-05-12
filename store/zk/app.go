package zk

import (
	"github.com/bbklab/swan-ng/types"
)

//
// app CRUD
//

// CreateApp ...
func (s *Store) CreateApp(app *types.App) error {
	bs, err := encode(app)
	if err != nil {
		return err
	}

	path := keyApp + "/" + app.ID
	return s.createAll(path, bs)
}

// UpdateApp ...
func (s *Store) UpdateApp(app *types.App) error {
	bs, err := encode(app)
	if err != nil {
		return err
	}

	path := keyApp + "/" + app.ID
	return s.create(path, bs)
}

// GetApp ...
func (s *Store) GetApp(id string) (*types.App, error) {
	bs, err := s.get(keyApp + "/" + id)
	if err != nil {
		return nil, err
	}

	app := new(types.App)
	if err := decode(bs, &app); err != nil {
		return nil, err
	}

	return app, nil
}

// ListApps ...
func (s *Store) ListApps() ([]*types.App, error) {
	nodes, err := s.list(keyApp)
	if err != nil {
		return nil, err
	}

	ret := make([]*types.App, 0, len(nodes))
	for _, node := range nodes {
		bs, err := s.get(keyApp + "/" + node)
		if err != nil {
			return nil, err
		}

		app := new(types.App)
		if err := decode(bs, &app); err != nil {
			return nil, err
		}

		ret = append(ret, app)
	}

	return ret, nil
}

// DeleteApp ...
func (s *Store) DeleteApp(id string) error {
	return s.del(keyApp + "/" + id)
}

//
// TODO app's version
//

// CreateVersion ...
func (s *Store) CreateVersion(aid string, ver *types.Version) error {
	return nil
}

// GetVersion ...
func (s *Store) GetVersion(aid, vid string) (*types.Version, error) {
	return nil, nil
}

// ListVersions ...
func (s *Store) ListVersions(aid string) ([]*types.Version, error) {
	return nil, nil
}

//
// TODO app's tasks
//

// ListTasks ...
func (s *Store) ListTasks(aid string) ([]*types.Task, error) {
	return nil, nil
}

// UpdateTask ...
func (s *Store) UpdateTask(aid string, t *types.Task) error {
	return nil
}

// GetTaskHistories ...
func (s *Store) GetTaskHistories(aid, tid string) ([]*types.Task, error) {
	return nil, nil
}
