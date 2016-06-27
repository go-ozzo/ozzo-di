// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package di_test

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
	Bar `inject:"true"`
}

func Example() {
	// creating a DI container
	c := di.NewContainer()

	// register a Foo instance as the Bar interface type
	c.RegisterAs(&Foo{"hello"}, di.InterfaceOf((*Bar)(nil)))

	// &Foo{"hello"} will be injected as the Bar parameter for test()
	c.Call(test)

	// create a MyBar object and inject its Bar field
	bar := c.Make(reflect.TypeOf(&MyBar{})).(Bar)
	fmt.Println(bar.String() + "2")

	// Output:
	// hello
	// hello2
}
