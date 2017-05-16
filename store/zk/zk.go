package zk

import (
	"encoding/json"
	"net/url"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	keyApp      = "/app"      // single app
	keyInstance = "/instance" // compose instance (group apps)
)

// Store represents zk store
type Store struct {
	url  *url.URL
	conn *zk.Conn
	acl  []zk.ACL
}

// New ...
// TODO test if need reconnect logic
func New(url *url.URL) (*Store, error) {
	s := &Store{
		url: url,
		acl: zk.WorldACL(zk.PermAll),
	}

	if err := s.initConnection(); err != nil {
		return nil, err
	}

	// create base keys nodes
	for _, node := range []string{keyApp, keyInstance} {
		if err := s.createAll(node, nil); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Store) initConnection() error {
	hosts := strings.Split(s.url.Host, ",")
	conn, connCh, err := zk.Connect(hosts, 5*time.Second)
	if err != nil {
		return err
	}

	// waiting for zookeeper to be connected.
	for event := range connCh {
		if event.State == zk.StateConnected {
			log.Info("connected to zookeeper succeed.")
			break
		}
	}

	s.conn = conn

	return nil
}

// with the prefix `s.url.Path` and clean the path
func (s *Store) clean(p string) string {
	if !strings.HasPrefix(p, s.url.Path) {
		p = s.url.Path + "/" + p
	}
	return path.Clean(p)
}

func (s *Store) get(path string) (data []byte, err error) {
	data, _, err = s.conn.Get(s.clean(path))
	return
}

func (s *Store) del(path string) error {
	exist, err := s.exist(path)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}
	return s.conn.Delete(s.clean(path), -1)
}

func (s *Store) list(path string) (children []string, err error) {
	children, _, err = s.conn.Children(s.clean(path))
	return
}

func (s *Store) exist(path string) (exist bool, err error) {
	exist, _, err = s.conn.Exists(s.clean(path))
	return
}

func (s *Store) createAll(path string, data []byte) error {
	path = s.clean(path)

	var (
		fields = strings.Split(path, "/")
		node   = "/"
	)

	// all of dir node
	for i, v := range fields[1:] {
		node += v
		if i >= len(fields[1:])-1 {
			break // the end node
		}
		err := s.create(node, nil)
		if err != nil {
			log.Errorf("create node: %s error: %v", node, err)
			return err
		}
		node += "/"
	}

	// the end data node
	return s.create(node, data)
}

func (s *Store) create(path string, data []byte) error {
	path = s.clean(path)

	exist, err := s.exist(path)
	if err != nil {
		return err
	}

	if exist {
		_, err = s.conn.Set(path, data, -1)
	} else {
		_, err = s.conn.Create(path, data, 0, s.acl)
	}

	return err
}

// encode & decode is just short-hands for json Marshal/Unmarshal
func encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}
func decode(bs []byte, v interface{}) error {
	return json.Unmarshal(bs, v)
}
