package admon

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/flect"
)

type Portal struct {
	options  Options
	registry []*Resource
}

func NewPortal(options Options) *Portal {
	return &Portal{
		options: options,
	}
}

func (p *Portal) registryMW(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Set("resources", p.registry)

		return next(c)
	}
}

// Add allows to add an admin resource with passed options for it.
func (p *Portal) Resource(model interface{}) *Resource {
	opts := ResourceOptions{}

	fieldr := NewFielder(model, opts.Fields)
	fieldr.portal = p

	resource := NewResource(model, fieldr)
	resource.Portal = p
	resource.Options = opts

	p.registry = append(p.registry, resource)
	return resource
}

// MountTo mounts the actual routes into the passed app, it will consider
// previously set Prefix and settings when adding.
func (p *Portal) MountTo(app *buffalo.App) {
	adminGroup := app.Group(p.options.Prefix)

	//TODO: add middlewares
	//TODO: security
	adminGroup.Use(p.registryMW)

	adminGroup.GET("/", func(c buffalo.Context) error {
		return c.Render(200, renderEngine.HTML("dashboard.html"))
	})

	for idx, res := range p.registry {
		ident := flect.New(structs.New(res.model).Name())

		base := res.Options.Path
		if base == "" {
			base = ident.Pluralize().Underscore().String()
		}

		paths := Paths{
			Prefix:   p.options.Prefix,
			BasePath: base,
		}

		res.Paths = paths
		p.registry[idx] = res

		//add the routes to the app.
		r := adminGroup.Group(base)
		r.GET("/", res.list)
		r.GET("/new", res.new)
		r.POST("/", res.create)
		r.GET("/export.{format}", res.export)

		r.GET(fmt.Sprintf("/{%v_id}", ident.Singularize().Underscore().String()), res.show)
		r.GET(fmt.Sprintf("/{%v_id}/edit", ident.Singularize().Underscore().String()), res.edit)
		r.PUT(fmt.Sprintf("/{%v_id}", ident.Singularize().Underscore().String()), res.update)
		r.DELETE(fmt.Sprintf("/{%v_id}", ident.Singularize().Underscore().String()), res.destroy)
	}

	as := NewAssetsServer(assetsBox, "/assets")
	as.AddHelpersTo(renderEngine)
	as.MountTo(adminGroup)
}
