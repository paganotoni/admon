package admon

import (
	"reflect"
	"time"

	"github.com/fatih/structs"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form/bootstrap"
)

type Fielder struct {
	Fields []*structs.Field

	fieldOptions FieldOptionsSet
}

func NewFielder(model interface{}, opts Options) Fielder {
	//TODO: handle options
	reflected := structs.New(model)

	return Fielder{
		Fields:       reflected.Fields(),
		fieldOptions: opts.Fields,
	}
}

func (fr Fielder) ValueFor(element interface{}, field *structs.Field, tx *pop.Connection) interface{} {

	raw := structs.New(element).Field(field.Name()).Value()

	switch v := raw.(type) {
	case time.Time:
		return v.Format(DateFormat)
	}

	for _, vr := range fr.fieldOptions {
		if vr.Name != field.Name() || vr.Renderer == nil {
			continue
		}

		tag, _ := vr.Renderer(raw, tx)
		return tag
	}

	return raw
}

func (fr Fielder) FormFields() []*structs.Field {
	result := []*structs.Field{}
	for _, f := range fr.Fields {
		//TODO: allow this to be customizable
		if f.Name() == "CreatedAt" || f.Name() == "UpdatedAt" || f.Name() == "ID" {
			continue
		}

		result = append(result, f)
	}

	return result
}

func (fr Fielder) TableFields() []*structs.Field {
	result := []*structs.Field{}

	for _, tc := range fr.fieldOptions {
		for _, f := range fr.Fields {
			if f.Name() == tc.Name {
				result = append(result, f)
			}
		}
	}

	if len(result) > 0 {
		return result
	}

	for _, f := range fr.Fields {
		if f.Name() == "ID" {
			continue
		}

		result = append(result, f)
	}

	return result
}

func (fr Fielder) SearchableFields() []*structs.Field {
	result := []*structs.Field{}

	for _, f := range fr.Fields {
		if f.Tag("db") == "-" || f.Tag("db") == "" {
			continue
		}

		if f.Kind() != reflect.String {
			continue
		}

		result = append(result, f)
	}

	return result
}

func (fr Fielder) LabelFor(field *structs.Field) string {
	label := field.Name()

	for _, fo := range fr.fieldOptions {
		if fo.Name != field.Name() {
			continue
		}

		if fo.Label == "" {
			break
		}

		label = fo.Label
	}

	return label
}

func (fr Fielder) FieldFor(element interface{}, field *structs.Field, form *bootstrap.FormFor, tx *pop.Connection) *tags.Tag {
	var opts FieldOptions
	for _, fo := range fr.fieldOptions {
		if fo.Name != field.Name() {
			continue
		}

		opts = fo
		break
	}

	switch opts.Input.Type {
	//TODO: add other types of fields
	case InputTypeSelect:
		//TODO: handle errors
		options, _ := opts.Input.SelectOptionsBuilder(tx)
		return form.SelectTag(field.Name(), tags.Options{"options": options, "hide_label": true})
	}

	return form.InputTag(field.Name(), tags.Options{"bootstrap": map[string]interface{}{"form-group-class": ""}, "hide_label": true})
}

func (fr Fielder) TableHeaderNameFor(field *structs.Field) string {
	return fr.LabelFor(field)
}

func (fr Fielder) ColumnNameFor(fieldName string) string {

	for _, f := range fr.Fields {
		if f.Name() != fieldName {
			continue
		}

		//TODO: handle cases like "-"
		return f.Tag("db")
	}

	return ""
}
