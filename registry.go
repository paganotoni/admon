package admon

import "github.com/paganotoni/admon/options"

var registry []*ResourceEntry

type ResourceEntry struct {
	Resource Resource
	Model    interface{}

	Options Options
}

func (reg *ResourceEntry) WithOptions(opts Options) *ResourceEntry {
	reg.Options = opts
	fieldr := NewFielder(reg.Model, opts)
	reg.Resource = NewResource(reg.Model, fieldr)
	return reg
}
