# ozzo-di

[![GoDoc](https://godoc.org/github.com/go-ozzo/ozzo-di?status.png)](http://godoc.org/github.com/go-ozzo/ozzo-di)
[![Build Status](https://travis-ci.org/go-ozzo/ozzo-di.svg?branch=master)](https://travis-ci.org/go-ozzo/ozzo-di)
[![Coverage](http://gocover.io/_badge/github.com/go-ozzo/ozzo-di)](http://gocover.io/github.com/go-ozzo/ozzo-di)

ozzo-di is a dependency injection (DI) container in Go language. It has the following features:

* DI via concrete types, interfaces, and provider functions
* DI of function parameter values and struct fields
* Creating and injecting new objects
* Hierarchical DI containers

## Requirements

Go 1.2 or above.

## Installation

Run the following command to install the package:

```
go get github.com/go-ozzo/ozzo-di
```

## Getting Started

The following code snippet shows how you can use the DI container.

```go
package main

import (
	"fmt"
	"reflect"
	"github.com/go-ozzo/ozzo-di"
)

type Bar interface {
    String() string
}

func test(bar Bar) {
    fmt.Println(bar.String())
}

type Foo struct {
    s string
}

func (f *Foo) String() string {
    return f.s
}

type MyBar struct {
    Bar `inject`
}

func main() {
    // creating a DI container
	c := di.NewContainer()

    // register a Foo instance as the Bar interface type
    c.RegisterAs(&Foo{"hello"}, di.InterfaceOf((*Bar)(nil)))

    // &Foo{"hello"} will be injected as the Bar parameter for test()
    c.Call(test)
    // Output:
    // hello

    // create a MyBar object and inject its Bar field
    bar := c.Make(reflect.TypeOf(&MyBar{})).(Bar)
    fmt.Println(bar.String())
    // Output:
    // hello
}
```


## Type Registration

`di.Container` is a DI container that relies on types to determine what values should be used for
injection. In order for this to happen, you usually should register the types that need DI support.
`di.Container` supports three kinds of type registration, as shown in the following code snippet:

```go
c := di.NewContainer()

// 1. registering a concrete type:

// &Foo{"hello"} is registered as the corresponding concrete type (*Foo)
c.Register(&Foo{"hello"})


// 2. registering an interface:

// &Foo{"hello"} is registered as the Bar interface
c.RegisterAs(&Foo{"hello"}, di.InterfaceOf((*Bar)(nil)))
// concrete type (*Foo) is registered as the Bar interface
c.RegisterAs(reflect.TypeOf(&Foo{}), di.InterfaceOf((*Bar)(nil)))


// 3. registering a provider:

// a provider function is registered as the Bar interface.
// The provider function will be called when injecting Bar.
c.RegisterProvider(func(di.Container) interface{} {
    return &Foo{"hello"}
}, di.InterfaceOf((*Bar)(nil)), true)
```

> Tip: To specify an interface type during registration, use the helper
> function `di.InterfaceOf((*InterfaceName)(nil))`.
> For concrete types, use go reflection function `reflect.TypeOf(TypeName{})`.


## Value Injection

`di.Container` supports three types of value injection, as shown in the following code snippet:

```go
// ...following the previous registration example...

type Composite struct {
    Bar `inject`
}

// 1. struct field injection:

// Exported struct field tagged with `inject` and anonymous fields will be injected with values.
// Here Composite.Bar will be injected with &Foo{"hello"}
composite := &Composite{}
c.Inject(composite)


// 2. function parameter injection:

// Function parameters will be injected with values according to their types.
// Here bar will be injected with &Foo{"hello"}
func test(bar Bar) {
    fmt.Println(bar.String())
}
c.Call(test)


// 3. making new instances:
// New struct instances can be created with their fields injected.
// Or a singleton instance may be returned.

foo := c.Make(reflect.TypeOf(&Foo{})).(*Foo)          // returns the singleton &Foo{"hello"}
bar := c.Make(di.InterfaceOf((*Bar)(nil))).(*Bar)     // returns the singleton &Foo{"hello"}

// returns a new Composite instance with Bar injected with the singleton &Foo{"hello"}
composite := c.Make(reflect.TypeOf(&Composite{})).(*Composite)
```

When injecting a previously registered type, if a value is already registered as that type, the value itself
will be used for injection.

If a provider is registered as a type, the provider will be called whose result will be used for injection.
While registering a provider, you may use the third parameter for `Container.RegisterProvider()` to indicate
whether the provider should be called every time the injection is needed or only the first time. If the
latter, the provider will only be called once and the same return result will be used for injection of
the corresponding registered type.

When injecting a value for a type `T` that has not been registered, the following strategy will be taken:

* If `*T` has been registered, the corresponding value will be dereferenced and returned;
* If `T` is a pointer type of `P`, the pointer to the value injected for `P` will be returned;
* If `T` is a struct type, a new instance will be created and its fields will be injected;
* If `T` is Slice, Map, or Chan, a new instance will be created and initialized;
* For all other cases, a zero value will be returned.


## Credits

ozzo-di has referenced the implementation of [codegansta/inject](https://github.com/codegangsta/inject/).
