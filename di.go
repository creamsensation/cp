package cp

import (
	"reflect"
)

type Provider interface {
	Provide(control Control) Provider
}

type Dependency struct {
	provider  Provider
	name      string
	singleton bool
}

const (
	diProvideMethodName = "Provide"
)

func Register(provider Provider) *Dependency {
	return &Dependency{
		provider: provider,
	}
}

func (d *Dependency) Name(name string) *Dependency {
	d.name = name
	return d
}

func (d *Dependency) Singleton() *Dependency {
	d.singleton = true
	return d
}

func Provide[T any](c Control, name ...string) *T {
	return provide[T](c.(*control), nil, name...).(*T)
}

func provide[T any](c *control, depType reflect.Type, name ...string) any {
	var n string
	if len(name) == 0 {
		n = reflect.TypeOf((*T)(nil)).String()
	}
	if len(name) > 0 {
		n = name[0]
	}
	dep, ok := c.core.deps[n]
	if !ok && depType != nil && depType.Kind() == reflect.Ptr {
		provider := reflect.New(depType.Elem())
		if !provider.IsValid() ||
			provider.IsValid() && !provider.MethodByName(diProvideMethodName).IsValid() {
			return nil
		}
		return provider.Interface().(Provider).Provide(c)
	}
	if dep == nil || (!ok && depType == nil) {
		return nil
	}
	if dep.singleton {
		return dep.provider
	}
	return dep.provider.Provide(c)
}
