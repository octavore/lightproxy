package service

import "reflect"

// A Module registers lifecycle hooks via the Config parameter of the Init
// method it must implement. A Module may depend on other Modules by declaring
// them as fields (pointers to other Modules). See Config for more information
// about the lifecycle hooks and the order in which they are executed. Within
// a Service, Modules are singletons.
type Module interface {
	Init(*Config)
}

// getModuleName returns the name of the module via reflection
func getModuleName(m Module) string {
	moduleType := reflect.TypeOf(m)
	if moduleType.Kind() == reflect.Ptr {
		moduleType = moduleType.Elem()
	}
	return moduleType.PkgPath() + "." + moduleType.Name() // todo: allow aliases
}
