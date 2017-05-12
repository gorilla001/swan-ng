// Package api ...
// events_manager.go implements a http events broadcasting system
// which serving http events stream response to subscribed clients
package api

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/bbklab/swan-ng/types"
)

var (
	// global events manager
	eventMgr *eventsManager
)

type eventsManager struct {
	sync.RWMutex                         // protect m
	m            map[string]*eventClient // store of online event clients
	max          int                     // max nb of clients, avoid bomber
}

type eventClient struct {
	w io.Writer
	f http.Flusher
	n http.CloseNotifier

	wait chan struct{}
	recv chan string
}

func init() {
	eventMgr = &eventsManager{
		m:   make(map[string]*eventClient),
		max: 1024,
	}

	// timer ping
	go func() {
		str := fmt.Sprintf("event: ping\ndata: \"it is a ping\"\n\n")
		for range time.Tick(10 * time.Second) {
			for _, c := range eventMgr.clients() {
				c.recv <- str
			}
		}
	}()
}

func (em *eventsManager) clients() map[string]*eventClient {
	em.RLock()
	defer em.RUnlock()
	return em.m
}

// broadcast message to all event clients
func (em *eventsManager) broadCast(e *types.Event) error {
	for _, c := range em.clients() {
		c.recv <- e.Format()
	}
	return nil
}

// subscribe() add an event client
func (em *eventsManager) subscribe(remoteAddr string, w io.Writer) {
	c := &eventClient{
		w: w,
		f: w.(http.Flusher),
		n: w.(http.CloseNotifier),

		wait: make(chan struct{}),
		recv: make(chan string, 1024),
	}

	go func(em *eventsManager, c *eventClient, remoteAddr string) {
		defer em.evict(remoteAddr)
		for {
			select {
			case <-c.n.CloseNotify():
				return
			case msg := <-c.recv:
				if _, err := c.w.Write([]byte(msg)); err != nil {
					log.Errorf("write event message to client [%s] error: [%v]", remoteAddr, err)
					return
				}
				c.f.Flush()
			}
		}
	}(em, c, remoteAddr)

	em.Lock()
	em.m[remoteAddr] = c
	em.Unlock()
}

func (em *eventsManager) wait(remoteAddr string) {
	em.RLock()
	c, ok := em.m[remoteAddr]
	em.RUnlock()
	if !ok {
		return
	}
	<-c.wait
}

func (em *eventsManager) evict(remoteAddr string) {
	log.Debug("evict event listener ", remoteAddr)
	em.Lock()
	defer em.Unlock()
	c, ok := em.m[remoteAddr]
	if !ok {
		return
	}
	close(c.wait)
	delete(em.m, remoteAddr)
}

func (em *eventsManager) full() bool {
	return em.size() >= int(em.max)
}

func (em *eventsManager) size() int {
	em.RLock()
	defer em.RUnlock()
	return len(em.m)
}
