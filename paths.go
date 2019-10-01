package admon

import (
	"fmt"
	"net/url"
	"path"
)

type Paths struct {
	Prefix   string
	BasePath string
}

func (a Paths) List() string {
	return a.Join(a.Prefix, a.BasePath)
}

func (a Paths) Show(id interface{}) string {
	return a.Join(a.Prefix, a.BasePath, fmt.Sprintf("%v", id))
}

func (a Paths) New() string {
	return a.Join(a.Prefix, a.BasePath, "new")
}

func (a Paths) Create() string {
	return a.Join(a.Prefix, a.BasePath)
}

func (a Paths) Edit(id interface{}) string {
	return a.Join(a.Prefix, a.BasePath, fmt.Sprintf("%v", id), "edit")
}

func (a Paths) Update(id interface{}) string {
	return a.Join(a.Prefix, a.BasePath, fmt.Sprintf("%v", id))
}

func (a Paths) Delete(id interface{}) string {
	return a.Join(a.Prefix, a.BasePath, fmt.Sprintf("%v", id))
}

func (a Paths) Export(format string) string {
	return a.Join(a.Prefix, a.BasePath, fmt.Sprintf("export.%v", format))
}

func (a Paths) Join(paths ...string) string {
	u, _ := url.Parse("/")
	u.Path = path.Join(paths...)
	return u.String()
}
