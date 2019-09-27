package admon

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/flect"
)

// Prefix is the prefix of the group where the admin
// resources will be mounted
var Prefix = "/admin"

// DateFormat that will be used for date fields and texts
var DateFormat = "01-02-2006"

// resources holds resources that we will use for admin routes generation.
var resources []*ResourceRegistry

type ResourceRegistry struct {
	Model    interface{}
	Resource Resource
	Options  Options
}

func (reg *ResourceRegistry) WithOptions(opts Options) *ResourceRegistry {
	reg.Options = opts
	fieldr := NewFielder(reg.Model, opts)
	reg.Resource = NewResource(reg.Model, fieldr)
	return reg
}

// Register allows to add an admin resource with passed options for it.
func Register(model interface{}) *ResourceRegistry {
	fieldr := NewFielder(model, Options{})
	resource := NewResource(model, fieldr)
	registry := &ResourceRegistry{
		Model:    model,
		Resource: resource,
		Options:  Options{},
	}

	resources = append(resources, registry)
	return registry
}

// MountTo mounts the actual routes into the passed app, it will consider
// previously set Prefix and settings when adding.
func MountTo(app *buffalo.App) {
	adminGroup := app.Group(Prefix)

	//TODO: mount index
	//TODO: add middlewares

	adminGroup.GET("/", func(c buffalo.Context) error {
		return c.Render(200, renderEngine.HTML("dashboard.html"))
	})

	for idx, res := range resources {
		ident := flect.New(structs.New(res.Model).Name())

		base := res.Options.Path
		if base == "" {
			base = ident.Pluralize().Underscore().String()
		}

		paths := Paths{
			Prefix:   Prefix,
			BasePath: base,
		}

		res.Resource.Paths = paths
		resources[idx] = res

		//add the routes to the app.
		r := adminGroup.Group(base)
		r.GET("/", res.Resource.list)
		r.GET("/new", res.Resource.new)
		r.POST("/", res.Resource.create)

		r.GET(fmt.Sprintf("/{%v_id}", ident.Singularize().Underscore().String()), res.Resource.show)
		r.GET(fmt.Sprintf("/{%v_id}/edit", ident.Singularize().Underscore().String()), res.Resource.edit)
		r.PUT(fmt.Sprintf("/{%v_id}", ident.Singularize().Underscore().String()), res.Resource.update)
		r.DELETE(fmt.Sprintf("/{%v_id}", ident.Singularize().Underscore().String()), res.Resource.destroy)
	}
}
