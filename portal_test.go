package admon_test

import (
	"net/http"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/httptest"
	"github.com/paganotoni/admon"
	"github.com/stretchr/testify/require"
)

func Test_Portal_Mount(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	portal := admon.NewPortal(admon.Options{
		Prefix: "/admin",
	})

	portal.AddTo(app)

	inf, err := app.Routes().Lookup("adminDashboardPath")
	r.NoError(err)
	r.NotNil(inf)
}

func Test_Portal_Assets(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	portal := admon.NewPortal(admon.Options{
		Prefix: "/admin",
	})

	portal.AddTo(app)

	htest := httptest.New(app)
	res := htest.HTML("/admin/assets/manifest.json").Get()
	r.Equal(http.StatusOK, res.Code)
}

func Test_Portal_Dashboard(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	portal := admon.NewPortal(admon.Options{
		Prefix: "/admin",
	})

	portal.AddTo(app)

	htest := httptest.New(app)
	res := htest.HTML("/admin/dashboard").Get()
	r.Equal(http.StatusOK, res.Code)
}
