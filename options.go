package admon

import "github.com/gobuffalo/tags"

const (
	InputTypeText   = 1
	InputTypeSelect = 2
)

//Options is the main struct for configuration on a particular admin resource.
type Options struct {
	//Path is the path prefix used for the routes under this resource.
	Path string
	//Fields allows customization for the fields that will be shown in the resource UI.
	Fields FieldOptionsSet
}

type FieldOptions struct {
	Name     string
	Input    InputOptions
	Renderer func(interface{}) *tags.Tag
}

type FieldOptionsSet []FieldOptions

type InputOptions struct {
	Type          int
	SelectOptions interface{}
}
