package goexports

import "reflect"

var Symbols = map[string]map[string]reflect.Value{}

//go:generate goexports github.com/osraige/visualisations
//go:generate goexports github.com/osraige/visualisations/timeline
//go:generate goexports github.com/osraige/visualisations/gauge
//go:generate goexports github.com/osraige/visualisations/clock
