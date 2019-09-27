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
	assetsBox    = packr.New("admin:assets", "./assets")
	templatesBox = packr.New("admin:templates", "./templates")

	helpers = render.Helpers{
		"resources":      registry,
		"addActiveClass": addActiveClass,
	}

	renderEngine = render.New(render.Options{
		HTMLLayout:   "admon.plush.html",
		TemplatesBox: templatesBox,

		AssetsBox: assetsBox,
		Helpers:   helpers,
	})
)

func registry() []*ResourceRegistry {
	return resources
}

func addActiveClass(entry *ResourceRegistry, help hctx.HelperContext) string {

	info := help.Value("current_route").(buffalo.RouteInfo)
	ident := flect.New(structs.New(entry.Model).Name())
	resourcePrefix := entry.Resource.Paths.Join(Prefix, ident.Pluralize().Underscore().String())

	if strings.HasPrefix(info.Path, resourcePrefix) {
		return "active"
	}

	return ""
}

