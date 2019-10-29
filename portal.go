package admon

import (
	"path/filepath"

	"github.com/gobuffalo/buffalo"
)

type Portal struct {
	options Options
	// registry []*Resource
}

func NewPortal(options Options) *Portal {
	portal := &Portal{
		options: options,
	}

	return portal
}

func (p *Portal) AddTo(app *buffalo.App) {
	portalGroup := app.Group(p.options.Prefix)
	portalGroup.GET("dashboard", p.dashboardHandler)
	app.ServeFiles(filepath.ToSlash(filepath.Join(p.options.Prefix, "assets")), assets)
}

func (p *Portal) dashboardHandler(c buffalo.Context) error {

	return nil
}
