package memory

import (
	"errors"

	"github.com/bbklab/swan-ng/types"
)

var (
	ErrAppNotExists = errors.New("app does not exist")
)

// CreateApp creates new app.
func (s *Store) CreateApp(app *types.App) error {
	s.Lock()
	defer s.Unlock()

	s.apps[app.ID] = app

	return nil
}

// UpdateApp update app information.
func (s *Store) UpdateApp(app *types.App) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.apps[app.ID]; !exists {
		return ErrAppNotExists
	}

	s.apps[app.ID] = app

	return nil
}

// GetApp get specified app by app id.
func (s *Store) GetApp(id string) (*types.App, error) {
	s.RLock()
	defer s.RUnlock()

	app, exists := s.apps[id]
	if !exists {
		return nil, ErrAppNotExists
	}

	return app, nil
}

// ListApps list all apps.
func (s *Store) ListApps() ([]*types.App, error) {
	s.RLock()
	defer s.RUnlock()

	apps := make([]*types.App, 0, len(s.apps))

	for _, v := range s.apps {
		apps = append(apps, v)
	}

	return apps, nil
}

// DeleteApp delete app by app id.
func (s *Store) DeleteApp(id string) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.apps[id]; !exists {
		return ErrAppNotExists
	}

	delete(s.apps, id)

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
