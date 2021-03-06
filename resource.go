package admon

import (
	"fmt"
	"reflect"
	"strings"

	"encoding/json"
	"encoding/xml"

	"net/http"

	"github.com/fatih/structs"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	csv "github.com/gocarina/gocsv"
)

type Resource struct {
	model interface{}

	ParamKey      string
	TitlePlural   string
	TitleSingular string

	Options ResourceOptions

	Fielder Fielder
	Paths   Paths
	Portal  *Portal
}

func NewResource(model interface{}, fieldr Fielder) *Resource {
	ident := flect.New(structs.New(model).Name())

	return &Resource{
		model: model,

		ParamKey:      ident.Singularize().Underscore().String() + "_id",
		TitlePlural:   ident.Pluralize().Titleize().Capitalize().String(),
		TitleSingular: ident.Singularize().Titleize().String(),

		Fielder: fieldr,
	}
}

func (r *Resource) WithOptions(opts ResourceOptions) *Resource {
	r.Options = opts
	r.Fielder = NewFielder(r.model, opts.Fields)
	return r
}

func (r *Resource) IsRelatedWith(routeInfo buffalo.RouteInfo) bool {
	ident := flect.New(structs.New(r.model).Name())
	resourcePrefix := r.Paths.Join(r.Portal.options.Prefix, ident.Pluralize().Underscore().String())

	return strings.HasPrefix(routeInfo.Path, resourcePrefix)
}

func (r *Resource) sortOrderFrom(c buffalo.Context) string {
	field := c.Param("sortBy")
	order := c.Param("order")

	field = r.Fielder.ColumnNameFor(field)

	if field == "" {
		field = "created_at"
	}

	if order == "" {
		order = "desc"
	}

	return fmt.Sprintf("%v %v", field, order)
}

func (r *Resource) searchScope(c buffalo.Context) pop.ScopeFunc {
	fields := r.Fielder.SearchableFields()
	term := c.Param("term")

	return func(q *pop.Query) *pop.Query {

		if term == "" {
			return q
		}

		clause := []string{}
		for _, field := range fields {
			column := field.Tag("db")
			clause = append(clause, fmt.Sprintf("%v LIKE '%%%v%%'", column, term))
		}

		q.Where(strings.Join(clause, " OR "))

		return q
	}
}

func (r *Resource) element() reflect.Value {
	return reflect.New(reflect.TypeOf(r.model))
}

func (r *Resource) slice() reflect.Value {
	return reflect.New(reflect.SliceOf(reflect.TypeOf(r.model)))
}

func (r *Resource) identifierFor(element interface{}) interface{} {
	return structs.New(element).Field("ID").Value()
}

func (r *Resource) list(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	q := tx.PaginateFromParams(c.Params())
	elements := r.slice().Interface()
	err := q.Scope(r.searchScope(c)).Order(r.sortOrderFrom(c)).All(elements)
	if err != nil {
		return err
	}

	c.Set("pagination", q.Paginator)
	c.Set("elements", elements)
	c.Set("resource", r)

	return c.Render(200, renderEngine.HTML("resource/index.plush.html"))
}

func (r *Resource) show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	element := r.element().Interface()
	err := tx.Find(element, c.Param(r.ParamKey))

	if err != nil {
		return c.Error(404, err)
	}

	c.Set("resource", r)
	c.Set("element", element)

	return c.Render(200, renderEngine.HTML("resource/show.plush.html"))
}

func (r *Resource) new(c buffalo.Context) error {

	c.Set("resource", r)
	c.Set("element", r.element().Interface())

	return c.Render(200, renderEngine.HTML("resource/new.plush.html"))
}

func (r *Resource) create(c buffalo.Context) error {

	element := r.element().Interface()
	if err := c.Bind(&element); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	//TODO: What's happening in pop.
	v := reflect.ValueOf(element)

	verrs, err := tx.ValidateAndCreate(v.Interface())
	if err != nil {
		return errors.Wrap(err, "error validating and creating")
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("resource", r)
		c.Set("element", element)

		return c.Render(200, renderEngine.HTML("resource/new.plush.html"))
	}

	c.Flash().Add("success", fmt.Sprintf("%v Created", r.TitlePlural))
	return c.Redirect(http.StatusSeeOther, r.Paths.Show(r.identifierFor(element)))
}

func (r *Resource) edit(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	element := r.element().Interface()
	if err := tx.Find(element, c.Param(r.ParamKey)); err != nil {
		return c.Error(404, err)
	}

	c.Set("element", element)
	c.Set("resource", r)

	return c.Render(200, renderEngine.HTML("resource/edit.plush.html"))
}

func (r *Resource) update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	element := r.element().Interface()
	if err := tx.Find(element, c.Param(r.ParamKey)); err != nil {
		return c.Error(404, err)
	}

	if err := c.Bind(&element); err != nil {
		return err
	}

	//TODO: What's happening in pop, this line should not be needed.
	v := reflect.ValueOf(element)

	verrs, err := tx.ValidateAndUpdate(v.Interface())
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("element", element)
		c.Set("resource", r)

		return c.Render(200, renderEngine.HTML("resource/edit.plush.html"))
	}

	c.Flash().Add("success", fmt.Sprintf("%v Updated", r.TitlePlural))
	return c.Redirect(http.StatusSeeOther, r.Paths.Show(r.identifierFor(element)))
}

func (r *Resource) destroy(c buffalo.Context) error {

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	element := r.element().Interface()

	if err := tx.Find(element, c.Param(r.ParamKey)); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(element); err != nil {
		return err
	}

	c.Flash().Add("success", fmt.Sprintf("%v Deleted", r.TitlePlural))
	return c.Redirect(http.StatusSeeOther, r.Paths.List())
}

func (r *Resource) export(c buffalo.Context) error {

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	elements := r.slice().Interface()
	err := tx.Scope(r.searchScope(c)).Order(r.sortOrderFrom(c)).All(elements)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(elements)
	format := c.Param("format")

	switch format {
	case "xml":
		data, _ = xml.MarshalIndent(elements, "  ", "    ")
	case "csv":
		data, _ = csv.MarshalBytes(elements)
	}

	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=export.%v", format))
	_, err = c.Response().Write(data)
	return err
}
