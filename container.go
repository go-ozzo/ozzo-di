// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package di implements a dependency injection (DI) container.
package di

import (
	"fmt"
	"reflect"
)

const injectTag = "inject"

// Container is a dependency injection (DI) container based on type mapping.
//
// Using Container involves two steps. First, register values, types, or providers with the types
// that should allow DI. Second, use one of the DI methods to achieve DI. For example,
//
//   import (
//       "reflect"
//       "github.com/go-ozzo/ozzo-di"
//   )
//
//   c := di.NewContainer()
//
//   // Step 1: register values, types, providers
//
//   // register the value Foo{"abc"} as type Foo
//   c.Register(Foo{"abc"})
//   // register the value Foo{"abc"} as the interface Bar
//   c.RegisterAs(Foo{"abc"}, di.InterfaceOf((*Bar)(nil)))
//   // register the struct Foo as the interface Bar
//   c.RegisterAs(reflect.TypeOf(Foo{}), di.InterfaceOf((*Bar)(nil)))
//   // register a provider that returns a shared value as the interface Bar
//   c.RegisterProvider(func(Container) reflect.Value {
//       return reflect.ValueOf(&Foo{"xyz"})
//   }, di.InterfaceOf((*Bar)(nil)), true)
//
//   // Step 2: dependency injection
//
//   // use `inject` tag to indicate which fields can be injected
//   type Tee struct {
//           Foo `inject`
//       bar Bar `inject`
//   }
//   // inject the fields of a struct
//   t := Tee{}
//   c.Inject(&t)
//   // inject function parameters
//   c.Call(func(bar Bar, foo Foo) {...})
//   // build a value of the specified type with injection
//   t2 := c.Build(reflect.TypeOf(&Tee{})).(*Tee)
//
// Note that when building an unregistered type, zero value will be returned. If the type is a struct,
// the zero value will be further injected by Inject() for those fields tagged with "inject".
type Container interface {
	// ParentContainer returns the parent container, if any.
	ParentContainer() Container
	// SetParentContainer sets the parent container.
	SetParentContainer(Container)

	// HasRegistered returns a value indicating whether the specified type has been registered before.
	HasRegistered(reflect.Type) bool
	// Unregister removes the specified type registration from the container
	Unregister(reflect.Type)

	// Register registers the specified value and associates it with the type of the value.
	Register(interface{})
	// RegisterAs registers the specified value or type, and associates it with the specified type.
	// For example,
	//
	//   c := di.NewContainer()
	//   // register a Foo struct as the Bar interface
	//   c.RegisterAs(&Foo{"abc"}, di.InterfaceOf((*Bar)(nil)))
	//   // register the Foo type as the Bar interface
	//   c.RegisterAs(reflect.TypeOf(&Foo{}), di.InterfaceOf((*Bar)(nil)))
	RegisterAs(interface{}, reflect.Type)
	// RegisterProvider registers the provider and associates it with the specified type.
	// When injecting or making a value for the type, the provider will be called and
	// its return value will be used as the value of the requested type. If shared is true,
	// the provider will only be called once, and its return value will be kept and used for
	// every injection request.
	RegisterProvider(p Provider, t reflect.Type, shared bool)

	// Call calls the specified function/method by injecting all its parameters.
	// The function/method result is returned as a slice.
	Call(interface{}) []interface{}
	// Inject injects the exported fields tagged with "inject" of the given struct.
	// Note that the struct should be passed as a pointer, or the fields won't be injected.
	Inject(interface{})
	// Make returns an instance of the specified type. If the instance is a newly created struct, its fields
	// will be injected by calling Inject(). Note that Make does not always create a new instance. If the type
	// has been registered and is associated with a value, that value will be returned.
	Make(reflect.Type) interface{}
}

// Provider is a function for creating a new instance of the associated type.
type Provider func(Container) reflect.Value

type providerBinding struct {
	provider Provider
	shared   bool
}

type container struct {
	parent Container
	values map[reflect.Type]interface{}
}

// NewContainer creates a new Dependency Injection (DI) container.
func NewContainer() Container {
	return &container{values: make(map[reflect.Type]interface{})}
}

// InterfaceOf is a helper method for turning an interface pointer into an interface reflection type.
// It is often used when calling RegisterAs() or Make() where you may want to specify an interface type.
// For example,
//
//   c := di.NewContainer()
//   c.RegisterAs(Foo{}, di.InterfaceOf((*Bar)(nil)))
//   foo := di.Make(di.InterfaceOf((*Bar)(nil)))
func InterfaceOf(iface interface{}) reflect.Type {
	t := reflect.TypeOf(iface)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Interface {
		panic("iface must be of an interface type, e.g., (*Foo)(nil)")
	}

	return t
}

func (c *container) ParentContainer() Container {
	return c.parent
}

func (c *container) SetParentContainer(parent Container) {
	c.parent = parent
}

func (c *container) HasRegistered(t reflect.Type) bool {
	_, ok := c.values[t]
	return ok
}

func (c *container) Unregister(t reflect.Type) {
	delete(c.values, t)
}

func (c *container) Register(val interface{}) {
	c.values[reflect.TypeOf(val)] = reflect.ValueOf(val)
}

func (c *container) RegisterAs(val interface{}, t reflect.Type) {
	if vt, ok := val.(reflect.Type); ok {
		// val is a type
		if !vt.ConvertibleTo(t) {
			panic(fmt.Sprintf("%v cannot be converted to %v", vt, t))
		}
		c.values[t] = vt
		return
	}

	vt := reflect.TypeOf(val)
	if !vt.ConvertibleTo(t) {
		panic(fmt.Sprintf("%v cannot be converted to %v", vt, t))
	}
	c.values[t] = reflect.ValueOf(val)
}

func (c *container) RegisterProvider(p Provider, t reflect.Type, shared bool) {
	c.values[t] = providerBinding{p, shared}
}

func (c *container) Call(f interface{}) []interface{} {
	t := reflect.TypeOf(f)

	// will panic if t is not a func while calling NumIn()
	var in = make([]reflect.Value, t.NumIn())

	for i := 0; i < t.NumIn(); i++ {
		in[i] = c.build(t.In(i))
	}

	s := reflect.ValueOf(f).Call(in)

	r := make([]interface{}, 0)
	for _, rv := range s {
		r = append(r, rv.Interface())
	}
	return r
}

func (c *container) Make(t reflect.Type) interface{} {
	return c.build(t).Interface()
}

func (c *container) Inject(val interface{}) {
	c.inject(reflect.ValueOf(val))
}

func (c *container) build(t reflect.Type) reflect.Value {
	if val, ok := c.values[t]; ok {
		switch val.(type) {
		case reflect.Type: // type mapped to interface
			return c.build(val.(reflect.Type))
		case providerBinding: // type mapped to provider func
			fb := val.(providerBinding)
			if !fb.shared {
				return fb.provider(c)
			}
			v := fb.provider(c)
			c.values[t] = v
			return v
		default: // type mapped to instance value
			return val.(reflect.Value)
		}
	}

	// no mapping found, try parent container, if any
	if c.parent != nil {
		return c.parent.(*container).build(t)
	}

	// try the pointer version
	ptr := reflect.PtrTo(t)
	if _, ok := c.values[ptr]; ok {
		return c.build(ptr).Elem()
	}

	// build from scratch
	switch t.Kind() {
	case reflect.Struct:
		r := reflect.New(t)
		c.inject(r)
		return r.Elem()
	case reflect.Slice:
		return reflect.MakeSlice(t, 0, 0)
	case reflect.Map:
		return reflect.MakeMap(t)
	case reflect.Chan:
		return reflect.MakeChan(t, 0)
	case reflect.Ptr:
		if v := c.build(t.Elem()); v.CanAddr() {
			return v.Addr()
		}
	}
	return reflect.New(t).Elem()
}

func (c *container) inject(v reflect.Value) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		tf := t.Field(i)
		if f.CanSet() && tf.Tag.Get(injectTag) != "" {
			f.Set(c.build(f.Type()))
		}
	}
}
