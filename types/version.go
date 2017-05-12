package types

import (
	"io"
	"text/template"
)

var versionTemplate = ` Version:      {{.Version}}
 Git commit:   {{.GitCommit}}
 Go version:   {{.GoVersion}}
 Built:        {{.BuildTime}}
 OS/Arch:      {{.Os}}/{{.Arch}}
`

// Version is exported
type Version struct {
	Version   string
	GoVersion string
	GitCommit string
	BuildTime string
	Os        string
	Arch      string
}

// FormatTo is exported
func (v Version) FormatTo(w io.Writer) error {
	tmpl, _ := template.New("version").Parse(versionTemplate)
	return tmpl.Execute(w, v)
}
