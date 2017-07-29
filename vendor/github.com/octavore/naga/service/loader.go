package service

import (
	"reflect"
)

// Loads m and its dependencies in topological order using DFS.
// Populated the modules slice with dependencies of m and m in topological order.
// Keeps state in s.configs.
func (s *Service) load(m Module) {
	moduleName := getModuleName(m)
	config, ok := s.configs[moduleName]
	if ok {
		panic("cycle: " + moduleName)
	} else if config != nil {
		BootPrintln("skipping already initialized module", moduleName)
		return
	}

	s.configs[moduleName] = nil
	// maybe also store a context.Context? Or some other boot time config.
	config = &Config{parent: m, service: s}
	BootPrintln("[service] initializing", moduleName)
	m.Init(config)

	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() != reflect.Ptr {
			continue
		}

		// check if field is a module
		val := reflect.New(f.Type().Elem())
		dep, ok := val.Interface().(Module)
		if !ok {
			continue
		}

		// warn on unexported modules
		if !f.CanSet() {
			BootPrintln("[service] warning: unexported module in", moduleName, f.String())
			continue
		}

		// reuse already initialized values
		if n := s.configs[getModuleName(dep)]; n != nil {
			f.Set(reflect.ValueOf(n.parent))
			continue
		}

		// set the field and initialize the field
		f.Set(val)
		s.load(dep)
	}
	s.modules = append(s.modules, m)
	s.configs[moduleName] = config
}
