package kernel

import (
	"fmt"
	"reflect"
	"time"
)

type Event interface {
	Version() int
	At() time.Time
	By() string
}

type Publisher interface {
	Publish(Event) error
}

func eventTyper(e Event) string {
	t := reflect.TypeOf(e)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return fmt.Sprintf("%v.%v", t.PkgPath(), t.Name())
}

func eventCTOR(e Event) func() Event {
	t := reflect.TypeOf(e)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return func() Event {
		return reflect.New(t).Interface().(Event)
	}
}
