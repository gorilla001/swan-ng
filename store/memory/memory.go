package memory

import (
	"sync"

	"github.com/bbklab/swan-ng/types"
)

// Store represents memory store
type Store struct {
	sync.RWMutex

	apps map[string]*types.App
}

// New ...
func New() *Store {
	return &Store{
		apps: make(map[string]*types.App),
	}
}
