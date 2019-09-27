package admon

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/helpers/hctx"

	"github.com/gobuffalo/packr/v2"
)

var (
	Helpers = render.Helpers{
		"resources": func() []*ResourceRegistry {
			return resources
		},

		"addActiveClass": func(entry *ResourceRegistry, help hctx.HelperContext) string {

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
		TemplatesBox: packr.New("admin:templates", "./templates"),

		//TODO: Need to find a way that this package has its own css files.
		AssetsBox: packr.New("admin:assets", "./assets"),

		Helpers: Helpers,
	})
)
