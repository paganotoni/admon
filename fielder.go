package admon

import (
	"time"

	"github.com/fatih/structs"
	"github.com/gobuffalo/flect"
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

func (fr Fielder) ValueFor(element interface{}, field *structs.Field) interface{} {
	raw := structs.New(element).Field(field.Name()).Value()

	switch v := raw.(type) {
	case time.Time:
		return v.Format(DateFormat)
	}

	for _, vr := range fr.fieldOptions {
		if vr.Name == field.Name() && vr.Renderer != nil {
			return vr.Renderer(raw)
		}
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

func (fr Fielder) FieldFor(element interface{}, field *structs.Field, form *bootstrap.FormFor) *tags.Tag {
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
		return form.SelectTag(field.Name(), tags.Options{"options": opts.Input.SelectOptions})
	}

	return form.InputTag(field.Name(), tags.Options{})
}

func (fr Fielder) TableHeaderNameFor(field *structs.Field) string {
	return flect.Humanize(field.Name())
}

func (fr Fielder) LabelFor(field *structs.Field) string {
	return flect.Humanize(field.Name())
}
