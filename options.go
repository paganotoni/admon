package admon

import "github.com/gobuffalo/tags"

const (
	InputTypeText   = 1
	InputTypeSelect = 2
)

type Options struct {
	//Path is the path prefix used for the routes under this resource.
	Path string

	Fields FieldOptionsSet
}

type FieldOptionsSet []FieldOptions

type FieldOptions struct {
	Name     string
	Input    InputOptions
	Renderer func(interface{}) *tags.Tag
}

type InputOptions struct {
	Type          int
	SelectOptions interface{}
}

// BasePath returns the base path for the routes being added.
// This will default to the inflected plural of the passed model
// if it's not specified.
//
// p.e:
// Team will use /teams
// Game will use /games
// Person will use /people
// func (a Admin) BasePath() string {
// 	if a.Options.Path != "" {
// 		return a.Options.Path
// 	}

// 	return a.Pluralize().Underscore().String()
// }
