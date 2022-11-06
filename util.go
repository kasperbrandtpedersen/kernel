package kernel

import (
	"fmt"
	"reflect"
)

type Cloner interface {
	Clone() any
}

type reflectCloner struct {
	proto reflect.Type
}

func (c *reflectCloner) Clone() any {
	return reflect.New(c.proto).Interface()
}

func cloner(i any) Cloner {
	if c, ok := i.(Cloner); ok {
		return c
	}

	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return &reflectCloner{
		proto: t,
	}
}

func eventTyper(e Event) string {
	t := reflect.TypeOf(e)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return fmt.Sprintf("%v %v", t.PkgPath(), t.Name())
}
