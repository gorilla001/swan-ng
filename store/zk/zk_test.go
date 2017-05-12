package zk

import (
	"net/url"
	"testing"
)

func TestZK(t *testing.T) {
	var (
		url, _ = url.Parse("zk://bbklab.net:2181/swan")
		err    error
		path   string
		data   []byte
	)

	// new
	s, err := New(url)
	if err != nil {
		t.Fatal(err)
	}

	// createAll
	path = "/a/b/c/../data"
	err = s.createAll(path, []byte(`data: data...`))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("created path %s", path)
	data, err = s.get(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s: %s", path, string(data))

	// create
	path = "/d"
	err = s.create(path, []byte(`d: data...`))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("created path %s", path)

	data, err = s.get(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s: %s", path, string(data))

	// list
	path = "/"
	nodes, err := s.list(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("children of %s: %v", path, nodes)

	// remove
	err = s.del("/not-exists")
	if err != nil {
		t.Fatal(err)
	}
	for _, node := range nodes {
		err := s.del(node)
		if err != nil {
			t.Fatal(err)
			continue
		}
		t.Logf("delete %s succeed", node)
	}
}
