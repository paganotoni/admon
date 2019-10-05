package admon

import (
	"fmt"
	"net/url"
	"path"
)

type Paths struct {
	Prefix   string
	BasePath string

	model  interface{}
	portal *Portal
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

func (a Paths) Export(format string, params map[string]string) string {
	path := a.Join(a.Prefix, a.BasePath, fmt.Sprintf("export.%v", format))
	parsed, err := url.Parse(path)
	if err != nil {
		return path
	}

	q := parsed.Query()
	q.Set("sortBy", params["orderBy"])
	q.Set("order", params["order"])
	q.Set("term", params["term"])

	parsed.RawQuery = q.Encode()
	return parsed.String()
}

func (a Paths) Join(paths ...string) string {
	u, _ := url.Parse("/")
	u.Path = path.Join(paths...)
	return u.String()
}
