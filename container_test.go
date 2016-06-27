// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package di

import (
	"reflect"
	"testing"
)

type Bar interface {
	test(int) string
}

type Foo struct {
	a string
}

func (f *Foo) test(b int) string {
	return f.a
}

type Writer interface {
	Write(string) string
}

type ResponseWriter struct {
	Foo *Foo `inject:"true"`
	t   int
}

func (*ResponseWriter) Write(s string) string {
	return s
}

type Request struct {
	Bar Bar `inject:"true"`
}

type Context struct {
	data string
}

type Controller struct {
	*Context  `inject:"true"`
	Response  Writer `inject:"true"`
	response  Writer `inject:"true"`
	Response2 Writer `inject:"true"`
	Response3 Writer

	Request Request `inject:"true"`
	action  string
}

func TestRegister(t *testing.T) {
	var foo1 Foo
	foo2 := Foo{}
	foo3 := &Foo{}
	foo4 := make([]Foo, 0)
	var foo5 [10]Foo
	var foo6 chan Foo
	var foo7 map[string]Foo
	var bar1 Bar = &Foo{}
	bar2 := (*Bar)(nil)

	tests := []struct {
		value  interface{}
		result reflect.Type
	}{
		{100, reflect.TypeOf(int(0))},
		{"abc", reflect.TypeOf(string(""))},
		{true, reflect.TypeOf(bool(false))},
		{100, reflect.TypeOf(int(0))},
		{foo1, reflect.TypeOf(Foo{})},
		{foo2, reflect.TypeOf(Foo{})},
		{foo3, reflect.TypeOf(&Foo{})},
		{foo4, reflect.TypeOf(make([]Foo, 0))},
		{foo5, reflect.TypeOf([10]Foo{})},
		{foo6, reflect.TypeOf(make(chan Foo, 0))},
		{foo7, reflect.TypeOf(make(map[string]Foo))},
		{bar1, reflect.TypeOf(&Foo{})},
		{bar2, reflect.TypeOf((*Bar)(nil))},
	}

	for _, test := range tests {
		c := NewContainer()
		c.Register(test.value)
		if !c.HasRegistered(test.result) {
			t.Errorf("Register(%v) failed, expected %v is registered", test.value, test.result)
		}
	}
}

func TestRegisterInterface(t *testing.T) {
	var foo Foo
	var bar *Bar = (*Bar)(nil)

	c := NewContainer()
	barType := InterfaceOf(bar)
	c.RegisterAs(&foo, barType)
	if !c.HasRegistered(barType) {
		t.Errorf("RegisterAs(%v, %v) failed, expected %v is registered", &foo, barType, barType)
	}
	c.Unregister(barType)
	if c.HasRegistered(barType) {
		t.Errorf("barType is still registered, expected unregistered")
	}

	fooType := reflect.TypeOf(&foo)
	if c.HasRegistered(fooType) {
		t.Errorf("RegisterAs(%v, %v) failed, expected %v is NOT registered", &foo, barType, fooType)
	}

	defer func() {
		if e := recover(); e == nil {
			t.Errorf("Expected a panic when registering an incompatible object")
		}
	}()
	writerType := InterfaceOf((*Writer)(nil))
	c.RegisterAs(&foo, writerType)
}

func TestRegisterInterface2(t *testing.T) {
	var foo Foo
	fooType := reflect.TypeOf(&foo)
	writerType := InterfaceOf((*Writer)(nil))
	c := NewContainer()
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("Expected a panic when registering an incompatible type")
		}
	}()
	c.RegisterAs(fooType, writerType)
}

func TestRegisterProvider(t *testing.T) {
	var f = func(Container) reflect.Value {
		return reflect.ValueOf(Foo{})
	}
	bar := (*Bar)(nil)
	funcType := reflect.TypeOf(f)
	barType := InterfaceOf(bar)

	c := NewContainer()
	c.RegisterProvider(f, barType, true)
	if !c.HasRegistered(barType) {
		t.Errorf("RegisterProvider(%v, %v) failed, expected %v is registered", f, barType, barType)
	}
	if c.HasRegistered(funcType) {
		t.Errorf("RegisterProvider(%v, %v) failed, expected %v is NOT registered", f, barType, funcType)
	}
}

func TestBuildBasic(t *testing.T) {
	tests := []struct {
		rt     reflect.Type
		result interface{}
	}{
		{reflect.TypeOf(int(0)), 0},
		{reflect.TypeOf(string("abc")), ""},
		{reflect.TypeOf(Foo{}), Foo{}},
	}

	c := NewContainer()
	for _, test := range tests {
		v := c.Make(test.rt)
		if v != test.result {
			t.Errorf("Make(%v): expected %#v, got %#v", test.rt, test.result, v)
		}
	}

	v := c.Make(reflect.TypeOf(&Foo{}))
	if reflect.ValueOf(v).IsNil() {
		t.Errorf("Make(%v) = nil, expected not nil", &Foo{})
	}
}

func TestBuildComplex(t *testing.T) {
	c := NewContainer()

	slice := c.Make(reflect.TypeOf([]int(nil))).([]int)
	if slice == nil || len(slice) != 0 {
		t.Errorf("Make([]int): expected []int{}, got %#v", slice)
	}

	m := c.Make(reflect.TypeOf(map[string]Foo{})).(map[string]Foo)
	if m == nil || len(m) != 0 {
		t.Errorf("Make(map[string]Foo): expected map[string]Foo{}, got %#v", slice)
	}

	ch := c.Make(reflect.TypeOf(make(chan Foo))).(chan Foo)
	if ch == nil || len(ch) != 0 {
		t.Errorf("Make(chan Foo): expected chan Foo{}, got %#v", ch)
	}
}

func TestBuildStruct(t *testing.T) {
	c := NewContainer()
	foo := Foo{"abc"}
	c.Register(foo)

	foo1 := c.Make(reflect.TypeOf(Foo{})).(Foo)
	if foo1.a != "abc" {
		t.Errorf("Make(Foo): expected Foo.a=%q, got %q", "abc", foo1.a)
	}

	// the registered version should remain unchanged (registering by value)
	foo1.a = "test"
	foo2 := c.Make(reflect.TypeOf(Foo{})).(Foo)
	if foo2.a != "abc" {
		t.Errorf("Make(Foo) after change: expected Foo.a=%q, got %q", "abc", foo2.a)
	}

	// registering by value doesn't mean registering the corresponding pointer
	foo3 := c.Make(reflect.TypeOf(&Foo{})).(*Foo)
	if foo3 != nil {
		t.Errorf("Make(Foo*) != nil, expected nil")
	}
}

func TestBuildStructPtr(t *testing.T) {
	c := NewContainer()
	foo := &Foo{"abc"}
	c.Register(foo)

	foo1 := c.Make(reflect.TypeOf(&Foo{})).(*Foo)
	if foo1.a != "abc" {
		t.Errorf("Make(*Foo): expected Foo.a=%q, got %q", "abc", foo1.a)
	}

	// the registered version should be unchanged (registering by pointer)
	foo1.a = "test"
	foo2 := c.Make(reflect.TypeOf(&Foo{})).(*Foo)
	if foo2.a != "test" {
		t.Errorf("Make(*Foo) after change: expected Foo.a=%q, got %q", "test", foo2.a)
	}

	// registering by value doesn't mean registering the corresponding pointer
	foo3 := c.Make(reflect.TypeOf(Foo{})).(Foo)
	if foo3.a != "test" {
		t.Errorf("Make(Foo).a = %q, expected %q", foo3.a, "abc")
	}
}

func TestBuildInterface(t *testing.T) {
	c := NewContainer()
	w := ResponseWriter{t: 1}
	writer := InterfaceOf((*Writer)(nil))
	c.RegisterAs(&w, writer)

	// make interface
	r := c.Make(writer).(*ResponseWriter)
	if r.t != 1 {
		t.Errorf("Make(Writer).t = %v, expected 1", r.t)
	}

	// make pointer
	r2 := c.Make(reflect.TypeOf(&w)).(*ResponseWriter)
	if r2 == nil {
		t.Errorf("Make(*ResponseWriter) = nil, expected not nil")
	}

	// making an unregistered concrete struct should result in zero struct
	r3 := c.Make(reflect.TypeOf(w)).(ResponseWriter)
	if r3.t != 0 {
		t.Errorf("Make(ResponseWriter).t = %v, expected 0", r3.t)
	}

	// making an unregistered interface should result in nil
	r4 := c.Make(InterfaceOf((*Bar)(nil)))
	if r4 != nil {
		t.Errorf("Make(Bar) = %#v, expected nil", r4)
	}
}

func TestBuildProvider(t *testing.T) {
	// shared
	c := NewContainer()
	c.RegisterProvider(func(Container) reflect.Value {
		return reflect.ValueOf(&Foo{"abc"})
	}, reflect.TypeOf(&Foo{}), true)
	foo1 := c.Make(reflect.TypeOf(&Foo{})).(*Foo)
	if foo1.a != "abc" {
		t.Errorf("Make(*Foo).a = %q, expected %q", foo1.a, "abc")
	}
	foo1.a = "xyz"
	foo2 := c.Make(reflect.TypeOf(&Foo{})).(*Foo)
	if foo1.a != foo2.a {
		t.Errorf("Changed shared Foo, Make(*Foo).a = %q, expected %q", foo2.a, foo1.a)
	}

	// not shared
	c = NewContainer()
	c.RegisterProvider(func(Container) reflect.Value {
		return reflect.ValueOf(&Foo{"abc"})
	}, reflect.TypeOf(&Foo{}), false)
	foo3 := c.Make(reflect.TypeOf(&Foo{})).(*Foo)
	if foo3.a != "abc" {
		t.Errorf("Make(*Foo).a = %q, expected %q", foo3.a, "abc")
	}
	foo3.a = "xyz"
	foo4 := c.Make(reflect.TypeOf(&Foo{})).(*Foo)
	if foo3.a == foo4.a {
		t.Errorf("Changed non-shared Foo, Make(*Foo).a = %q, expected %q", foo4.a, "abc")
	}
}

func TestBuildInject(t *testing.T) {
	barType := InterfaceOf((*Bar)(nil))
	writerType := InterfaceOf((*Writer)(nil))

	c := NewContainer()
	c.Register(&Context{"abc"})
	c.RegisterAs(&Foo{"123"}, barType)
	c.RegisterAs(reflect.TypeOf(&ResponseWriter{}), writerType)

	con := c.Make(reflect.TypeOf(Controller{})).(Controller)
	if con.action != "" {
		t.Errorf("Controller.action = %q, expected %q", con.action, "")
	}
	if con.data != "abc" {
		t.Errorf("Controller.data = %q, expected %q", con.data, "abc")
	}
	if con.Request.Bar.test(0) != "123" {
		t.Errorf("Controller.request.bar.test() = %q, expected %q", con.Request.Bar.test(0), "123")
	}
}

func TestInject(t *testing.T) {
	writerType := InterfaceOf((*Writer)(nil))

	c := NewContainer()
	c.Register(&Context{"abc"})
	c.RegisterAs(&ResponseWriter{t: 123}, writerType)
	c.RegisterAs(&Foo{"xyz"}, InterfaceOf((*Bar)(nil)))
	con := &Controller{}
	c.Inject(con)
	// anonymous field injection
	if con.data != "abc" {
		t.Errorf("Controller.data = %q, expected %q", con.data, "abc")
	}
	// valid injection
	if con.Response == nil {
		t.Errorf("Controller.Response should not be nil")
	}
	if con.Response.(*ResponseWriter).t != 123 {
		t.Errorf("Controller.Response.t = %q, expected %q", con.Response.(*ResponseWriter).t, 123)
	}
	// inject struct tag
	if con.Response2 == nil {
		t.Errorf("Controller.Response2 should not be nil")
	}
	// field not marked with "inject" should not be injected
	if con.Response3 != nil {
		t.Errorf("Controller.Response3 should be nil")
	}
	// unexported should not be injected
	if con.response != nil {
		t.Errorf("Controller.response should be nil")
	}

	// recursive injection
	if con.Request.Bar == nil {
		t.Errorf("Controller.Request.bar = nil, expected not nil")
	}
	if con.Request.Bar.(*Foo).a != "xyz" {
		t.Errorf("Controller.Request.bar.(*Foo).a = %q, expected %q", con.Request.Bar.(*Foo).a, "xyz")
	}

	// nothing should happen for injecting a non-struct variable
	c.Inject(1)
}

func TestCall(t *testing.T) {
	writerType := InterfaceOf((*Writer)(nil))

	c := NewContainer()
	w := ResponseWriter{t: 1}
	context := Context{"abc"}
	c.RegisterAs(&w, writerType)
	c.Register(&context)

	f := func(cc *Context, w Writer, s string) (string, string, string) {
		return cc.data, w.Write("test"), s
	}
	result := c.Call(f)
	if len(result) != 3 {
		t.Error("The return result should contain 3 values.")
	}
	if result[0] != "abc" || result[1] != "test" {
		t.Errorf("The return result was (%q, %q, %q), expected (%q, %q, %q)", result[0], result[1], result[2], "abc", "test", "")
	}
}

func TestParent(t *testing.T) {
	c := NewContainer()
	if c.ParentContainer() != nil {
		t.Error("ParentContainer() != nil, expected nil")
	}

	f := c.Make(reflect.TypeOf(Foo{})).(Foo)
	if f.a != "" {
		t.Error("Build(Foo) returns non-zero Foo struct, expected zero")
	}

	p := NewContainer()
	p.Register(Foo{"abc"})
	c.SetParentContainer(p)
	if c.ParentContainer() == nil {
		t.Error("ParentContainer() == nil, expected not nil")
	}
	f = c.Make(reflect.TypeOf(Foo{})).(Foo)
	if f.a != "abc" {
		t.Errorf("Build(Foo).a=%q, expected %q", f.a, "abc")
	}
}

func TestInterfaceOf(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("Expected a panic when calling InterfaceOf() with a non-interface")
		}
	}()
	InterfaceOf(reflect.TypeOf(Foo{}))
}
