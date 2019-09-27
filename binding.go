package admon

import (
	"net/http"
	"time"

	"github.com/go-playground/form"
	"github.com/gobuffalo/buffalo/binding"
	"github.com/pkg/errors"
)

var decoder = form.NewDecoder()

func init() {
	decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		//TODO: support multiple date formats
		return time.Parse("2006-01-02", vals[0])
	}, time.Time{})

	formBinding := func(req *http.Request, i interface{}) error {
		err := req.ParseForm()
		if err != nil {
			return errors.WithStack(err)
		}

		if err := decoder.Decode(i, req.Form); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}

	binding.Register("application/html", formBinding)
	binding.Register("text/html", formBinding)
	binding.Register("application/x-www-form-urlencoded", formBinding)
	binding.Register("html", formBinding)
}
