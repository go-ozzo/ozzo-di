# ozzo-di

[![GoDoc](https://godoc.org/github.com/go-ozzo/ozzo-di?status.png)](http://godoc.org/github.com/go-ozzo/ozzo-di)
[![Build Status](https://travis-ci.org/go-ozzo/ozzo-di.svg?branch=master)](https://travis-ci.org/go-ozzo/ozzo-di)
[![Coverage](http://gocover.io/_badge/github.com/go-ozzo/ozzo-di)](http://gocover.io/github.com/go-ozzo/ozzo-di)

## 其他语言

[简体中文](/docs/README-zh-CN.md)
[Русский](/docs/README-ru.md)

## 说明

ozzo-di is a dependency injection (DI) container in Go language. It has the following features:
ozzo-di 是一个使用 Go 语言实现的依赖注入（DI）容器。它包含以下功能：

* DI via concrete types, interfaces, and provider functions
* 支持通过具体的类型(concrete types)、接口(interfaces)、以及提供函数(provider functions)进行注入
* DI of function parameter values and struct fields
* 对函数参数以及结构字段 (struct fields) 的注入
* Creating and injecting new objects
* 新对象的创建与注入
* Hierarchical DI containers
* 分层的 DI 容器(Hierarchical DI containers)

## 需求

Go 1.2 或以上。

## 安装

请运行以下指令安装此包：

```
go get github.com/go-ozzo/ozzo-di
```

## 准备开始

The following code snippet shows how you can use the DI container.
以下代码片段展示了可以如何使用该 DI 容器

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
	// 新建一个 DI 容器
	c := di.NewContainer()

    // register a Foo instance as the Bar interface type
	// 向 Bar 接口类型注册一个名为 Foo 的实例对象
    c.RegisterAs(&Foo{"hello"}, di.InterfaceOf((*Bar)(nil)))

    // &Foo{"hello"} will be injected as the Bar parameter for test()
	// &Foo{"hello"} 会被作为 Bar 参数注入到 test() 方法
    c.Call(test)
    // 输出：
    // hello

    // 创建一个 MyBar 对象并注入其 Bar 字段
    bar := c.Make(reflect.TypeOf(&MyBar{})).(Bar)
    fmt.Println(bar.String())
    // 输出：
    // hello
}
```


## 类型注册

`di.Container` is a DI container that relies on types to determine what values should be used for
injection. In order for this to happen, you usually should register the types that need DI support.
`di.Container` supports three kinds of type registration, as shown in the following code snippet:
`di.Container` 是一个通过对象的类型来判断需要注入什么对象的 DI 容器。为了实现这一点，通常你需要先注册这些类型，才能使其
用于依赖注入。`di.Container` 支持 3 中不同的类型注册方式，如下所示：

```go
c := di.NewContainer()

// 1. 注册一个具体类型：

// &Foo{"hello"} 被注册为其对应的具体类型 (*Foo)
c.Register(&Foo{"hello"})


// 2. 注册一个接口：

// &Foo{"hello"} 被注册为 Bar 接口
c.RegisterAs(&Foo{"hello"}, di.InterfaceOf((*Bar)(nil)))
// concrete type (*Foo) is registered as the Bar interface
c.RegisterAs(reflect.TypeOf(&Foo{}), di.InterfaceOf((*Bar)(nil)))


// 3. 注册一个提供方法：

// 一个提供方法被注册为 Bar 接口
// 该方法会在需要注入 Bar 的时候被调用。
c.RegisterProvider(func(di.Container) interface{} {
    return &Foo{"hello"}
}, di.InterfaceOf((*Bar)(nil)), true)
```

> Tip: 要在类型注册时，指定一个接口类型你可以使用助手方法 `di.InterfaceOf((*InterfaceName)(nil))`。
> 要指定具体的类型，你可以使用 golang 的反射方法 `reflect.TypeOf(TypeName{})`.


## 值注入

`di.Container` supports three types of value injection, as shown in the following code snippet:
`di.Container` 支持三种值注册，他们分别列举在下面的代码片段中:

```go
// ...跟着之前类型注册的例子...

type Composite struct {
    Bar `inject`
}

// 1. struct field injection:
// 1. 结构字段的注入:

// Exported struct field tagged with `inject` and anonymous fields will be injected with values.
// 容器可以注入使用 `inject` 标记的公开结构字段，以及匿名字段。
// Here Composite.Bar will be injected with &Foo{"hello"}
composite := &Composite{}
c.Inject(composite)


// 2. function parameter injection:
// 2. 函数参数的注入:

// Function parameters will be injected with values according to their types.
// 函数的参数会依据其参数的类型被注入
// Here bar will be injected with &Foo{"hello"}
// 这里 bar 会被注入 &Foo{"hello"} 对象 （译者注，根据上面的例子，&Foo("hello") 是一个 Bar 对象）
func test(bar Bar) {
    fmt.Println(bar.String())
}
c.Call(test)


// 3. making new instances:
// 3. 初始化新的实例:
// New struct instances can be created with their fields injected.
// 新结构实例会在初始化的时候注入其子字段的值。
// Or a singleton instance may be returned.
// 或者返回一个单例对象。

foo := c.Make(reflect.TypeOf(&Foo{})).(*Foo)          // 返回单例 &Foo{"hello"}
bar := c.Make(di.InterfaceOf((*Bar)(nil))).(*Bar)     // 返回单例 &Foo{"hello"}

// returns a new Composite instance with Bar injected with the singleton &Foo{"hello"}
// 返回一个新的 Composite 对象，它的 Bar 字段会被注入一个单例 &Foo{"hello"} 对象
composite := c.Make(reflect.TypeOf(&Composite{})).(*Composite)
```

When injecting a previously registered type, if a value is already registered as that type, the value itself
will be used for injection.
当注入一个之前注册过的类型时，若某对象已经被注册给这个类型了，那么这个对象就会被用于注入。

If a provider is registered as a type, the provider will be called whose result will be used for injection.
若注册到某类型的是一个提供方法，那么这个提供方法被调用之后的返回值就会被用于注入。
While registering a provider, you may use the third parameter for `Container.RegisterProvider()` to indicate
whether the provider should be called every time the injection is needed or only the first time. 
当注册一个提供方法时，你可以选择使用 `Container.RegisterProvider()` 方法的第三个参数，来指明是在每次注册前都调用提供方法，还是只在第一次调用的时候。
If the latter, the provider will only be called once and the same return result will be used for injection of
the corresponding registered type.
如果是后者，那么提供方法只会被调用一次，而这之后每次该类型的注入都值会使用之前那个返回值。

When injecting a value for a type `T` that has not been registered, the following strategy will be taken:
当发现需要注入给类型 `T` 的值，之前并没有被注册过的话，那么容器会依次以如下策略进行判断：

* If `*T` has been registered, the corresponding value will be dereferenced and returned;
* 若 `*T` 被注册过，那么会返回其对应的值，并解除引用；
* If `T` is a pointer type of `P`, the pointer to the value injected for `P` will be returned;
* 若 `T` 是 `P` 的指针类型，那么指向注册给 `P` 类型的那个对象的指针会被返回
* If `T` is a struct type, a new instance will be created and its fields will be injected;
* 若 `T` 是一个结构类型，会创建一个新的实例，并且它的子字段都会被注入。
* If `T` is Slice, Map, or Chan, a new instance will be created and initialized;
* 如果 `T` 是一个 Slice, Map 或者 Chan, 会初始化并返回一个新实例。
* For all other cases, a zero value will be returned.
* 其他情况，会返回零值。


## 鸣谢

ozzo-di 参考了 [codegansta/inject](https://github.com/codegangsta/inject/) 的实现。
