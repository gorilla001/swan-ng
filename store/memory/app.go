package memory

import (
	"github.com/bbklab/swan-ng/types"
)

// CreateApp ...
func (s *Store) CreateApp(app *types.App) error {
	return nil
}

// UpdateApp ...
func (s *Store) UpdateApp(app *types.App) error {
	return nil
}

// GetApp ...
func (s *Store) GetApp(id string) (*types.App, error) {
	return nil, nil
}

// ListApps ...
func (s *Store) ListApps() ([]*types.App, error) {
	return nil, nil
}

// DeleteApp ...
func (s *Store) DeleteApp(id string) error {
	return nil
}

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
