package nav

import (
	. "github.com/snivilised/extendio/translate"
)

// InitFiltersHookFn is the default filter initialiser. This can be overridden or extended
// by the client if the need arises. To extend this behaviour rather than replace it,
// call this function from inside the custom function set on o.Hooks.Filter. To
// replace the default functionality, do note the following points:
// - the original client callback is defined as frame.client, this should be referred to
// from outside the custom function (ie in the closure) as is performed here in the default.
// This will allow the custom function to invoke the core callback as appropriate.
// - The filters defined here in extendio make use of some extended fields, so if the client
// needs to define a custom function that is compatible with the native filters, then make
// sure the DoExtend value is set to true in the options, otherwise a panic will occur due to the
// filter attempting to de-reference the Extension on the TraverseItem.
func InitFiltersHookFn(o *TraverseOptions, frame *navigationFrame) {

	if o.Store.FilterDefs != nil {
		if o.Store.FilterDefs.Current.Source != "" || o.Store.FilterDefs.Current.Custom != nil {
			o.useExtendHook()
			frame.filters.Current = NewCurrentFilter(&o.Store.FilterDefs.Current)
			frame.filters.Current.Validate()
			decorated := frame.client
			decorator := func(item *TraverseItem) *LocalisableError {
				if frame.filters.Current.IsMatch(item) {
					return decorated(item)
				}
				return nil
			}
			frame.decorate("init-current-filter 🎁", decorator)
		}

		if o.Store.FilterDefs.Children.Source != "" || o.Store.FilterDefs.Children.Custom != nil {
			o.useExtendHook()
			frame.filters.Compound = NewCompoundFilter(&o.Store.FilterDefs.Children)
			frame.filters.Compound.Validate()
		}
	}
}

func bootstrapFilter(o *TraverseOptions, frame *navigationFrame) {
	o.Hooks.InitFilters(o, frame)
}
