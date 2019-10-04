package admon

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/tags"
)

const (
	//InputTypeText is the text input
	InputTypeText     = 1
	InputTypeSelect   = 2
	InputTypeTextarea = 3

	InputTypeCheckbox    = 4
	InputTypeRadioButton = 5
)

//Options is the main struct for configuration on a particular admin resource.
type Options struct {
	//Path is the path prefix used for the routes under this resource.
	Path string
	//Fields allows customization for the fields that will be shown in the resource UI.
	Fields FieldOptionsSet
}

type FieldOptions struct {
	// Name of the field we're defining config. This is the way we know which
	// one is being configured. This is the same value at the struct level.
	Name string
	// Label that will be used in the table and in the form
	Label string

	// OnlyForm
	OnlyForm bool

	Input    InputOptions
	Renderer func(interface{}, *pop.Connection) (*tags.Tag, error)
}

type FieldOptionsSet []FieldOptions

type InputOptions struct {
	//Type of the field that will be used
	Type int

	//SelectOptionsBuilder receives the database connection and creates the options
	SelectOptionsBuilder func(tx *pop.Connection) (interface{}, error)
}
