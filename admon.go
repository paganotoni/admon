package admon

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/helpers/hctx"
	"github.com/gobuffalo/packr/v2"
)

var (
	// Prefix is the prefix of the group where the admin
	// resources will be mounted
	Prefix = "/admin"

	// DateFormat that will be used for date fields and texts
	DateFormat = "01-02-2006"

	assetsBox    = packr.New("admon:assets", "./web/public/assets")
	templatesBox = packr.New("admon:templates", "./web/templates")

	helpers = render.Helpers{
		"resources": func() []*ResourceEntry {
			return registry
		},
		"addActiveClass": func(entry *ResourceEntry, help hctx.HelperContext) string {

			info := help.Value("current_route").(buffalo.RouteInfo)
			ident := flect.New(structs.New(entry.Model).Name())
			resourcePrefix := entry.Resource.Paths.Join(Prefix, ident.Pluralize().Underscore().String())

			if strings.HasPrefix(info.Path, resourcePrefix) {
				return "active"
			}

			return ""
		},
	}

	renderEngine = render.New(render.Options{
		HTMLLayout:   "admon.plush.html",
		TemplatesBox: templatesBox,

		AssetsBox: assetsBox,
		Helpers:   helpers,
	})
)

// Register allows to add an admin resource with passed options for it.
func Register(model interface{}) *ResourceEntry {
	opts := Options{}
	fieldr := NewFielder(model, opts)
	resource := NewResource(model, fieldr)
	entry := &ResourceEntry{
		Model:    model,
		Resource: resource,
		Options:  opts,
	}

	registry = append(registry, entry)
	return entry
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

	for idx, res := range registry {
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
		registry[idx] = res

		//add the routes to the app.
		r := adminGroup.Group(base)
		r.GET("/", res.Resource.list)
		r.GET("/new", res.Resource.new)
		r.POST("/", res.Resource.create)

		r.GET("/export.{format}", res.Resource.export)

		r.GET(fmt.Sprintf("/{%v_id}", ident.Singularize().Underscore().String()), res.Resource.show)
		r.GET(fmt.Sprintf("/{%v_id}/edit", ident.Singularize().Underscore().String()), res.Resource.edit)
		r.PUT(fmt.Sprintf("/{%v_id}", ident.Singularize().Underscore().String()), res.Resource.update)
		r.DELETE(fmt.Sprintf("/{%v_id}", ident.Singularize().Underscore().String()), res.Resource.destroy)
	}

	as := NewAssetsServer(assetsBox, "/assets")
	as.AddHelpersTo(renderEngine)
	as.MountTo(adminGroup)
}
