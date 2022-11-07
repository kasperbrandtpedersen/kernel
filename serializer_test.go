package kernel_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/kasperbrandtpedersen/kernel"
)

type serializerEvent struct {
	EventVersion int
	EventAt      time.Time
	EventBy      string
}

func (e *serializerEvent) Version() int {
	return e.EventVersion
}

func (e *serializerEvent) At() time.Time {
	return e.EventAt
}

func (e *serializerEvent) By() string {
	return e.EventBy
}

func TestSerializer(t *testing.T) {
	serializer := kernel.NewJSONSerializer()

	serializer.Bind(&serializerEvent{})

	expected := &serializerEvent{
		EventVersion: 1,
		EventAt:      time.Now(),
		EventBy:      "Foo Bar",
	}
	rec, err := serializer.Serialize(expected)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := serializer.Deserialize(rec)

	if err != nil {
		t.Fatal(err)
	}

	if !equals(expected, actual) {
		t.Fatal()
	}
}

func BenchmarkDeserialize(b *testing.B) {
	serializer := kernel.NewJSONSerializer()

	serializer.Bind(&serializerEvent{})

	event := &serializerEvent{
		EventVersion: 1,
		EventAt:      time.Now(),
		EventBy:      "Foo Bar",
	}
	rec, err := serializer.Serialize(event)

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		serializer.Deserialize(rec)
	}
}

func BenchmarkSerialize(b *testing.B) {
	serializer := kernel.NewJSONSerializer()

	serializer.Bind(&serializerEvent{})

	event := &serializerEvent{
		EventVersion: 1,
		EventAt:      time.Now(),
		EventBy:      "Foo Bar",
	}

	for i := 0; i < b.N; i++ {
		serializer.Serialize(event)
	}
}

func equals(expected, actual kernel.Event) bool {
	if reflect.TypeOf(expected).Elem() != reflect.TypeOf(actual).Elem() {
		return false
	}

	if expected.Version() != actual.Version() {
		return false
	}

	if expected.By() != actual.By() {
		return false
	}

	return expected.At().Equal(actual.At())
}
