package admon

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/helpers/hctx"
	"github.com/gobuffalo/packr/v2"
)

var (
	// DefaultLabels allows to define labels for repeated fields like
	// ID, CreatedAt and UpdatedAt. This can be overriden if field
	// level label is defined.
	DefaultLabels = map[string]string{}

	assetsBox    = packr.New("admon:assets", "./web/public/assets")
	templatesBox = packr.New("admon:templates", "./web/templates")

	helpers = render.Helpers{
		"addActiveClass": func(resource *Resource, help hctx.HelperContext) string {

			route := help.Value("current_route").(buffalo.RouteInfo)
			if resource.IsRelatedWith(route) {
				return "active"
			}

			return ""
		},

		"pathWith": func(params map[string]interface{}, help hctx.HelperContext) template.HTML {
			request := help.Value("request").(*http.Request)
			path := request.URL

			q := path.Query()
			for name, value := range params {
				q.Set(name, fmt.Sprintf("%v", value))
			}

			path.RawQuery = q.Encode()

			return template.HTML(path.String())
		},
	}

	renderEngine = render.New(render.Options{
		HTMLLayout:   "admon.plush.html",
		TemplatesBox: templatesBox,

		AssetsBox: assetsBox,
		Helpers:   helpers,
	})
)
