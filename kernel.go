package kernel

import (
	"reflect"
	"sync"
)

// Missing features, maybe
// CLOCK, cron job-ish and/or do stuff at specific point in time. [need to have]
// QUEUES, in/outbox for error queues, maybe. Could ease async processing. [nice to have]
// GLOBAL, mechanisnm for global versioning across streams. [nice to have]

type Kernel struct {
	serializer  Serializer
	store       Store
	cmdRoutes   map[reflect.Type]func(Command) error
	eventRoutes map[reflect.Type][]func(Event) error
}

type Option func(*Kernel)

func New(options ...Option) *Kernel {
	k := &Kernel{
		serializer: &jsonSerializer{
			mu:    sync.RWMutex{},
			types: map[string]Cloner{},
		},
		store: &store{
			mu:      sync.RWMutex{},
			content: map[string]History{},
		},
		cmdRoutes:   map[reflect.Type]func(Command) error{},
		eventRoutes: map[reflect.Type][]func(Event) error{},
	}

	for _, opt := range options {
		opt(k)
	}

	return k
}

// TODO
// Maybe do async processing to avoid long running Dispatch/Publish.
// Dispatching or publihing may cause very deep Dispatch/Publish loop.

// TODO
// Maybe return something more relevent than just error.
// New state version, produced events or stuff like that.

func (k *Kernel) Dispatch(cmd Command) error {
	return k.cmdRoutes[reflect.TypeOf(cmd)](cmd)
}

func (k *Kernel) Publish(e Event) error {
	routes, ok := k.eventRoutes[reflect.TypeOf(e)]

	if !ok {
		return nil
	}

	var err error

	for _, r := range routes {
		err = r(e)
	}

	return err
}

func WithStore(s Store) Option {
	return func(k *Kernel) {
		k.store = s
	}
}

func WithSerializer(s Serializer) Option {
	return func(k *Kernel) {
		k.serializer = s
	}
}

// TODO
// Maybe use some form of Emitter over []Event returns.
// Return values of []Event gets pretty cumbersome to code and read.

func Decide[S State, CMD Command](decider func(s S, cmd CMD) []Event) Option {
	var proto CMD

	return func(k *Kernel) {
		dispatcher := func(cmd Command) error {
			repo := repository[S]{
				serializer: k.serializer,
				store:      k.store,
				publisher:  k,
			}
			stream := cmd.Stream()
			state, err := repo.Load(stream)

			if err != nil {
				return err
			}

			events := decider(state, cmd.(CMD))

			if err := repo.Save(stream, events); err != nil {
				return err
			}

			return nil
		}

		k.cmdRoutes[reflect.TypeOf(proto)] = dispatcher
	}
}

func Evolve[E Event](evolver func(e E, d Dispatcher) error) Option {
	var proto E

	return func(k *Kernel) {
		handler := func(event Event) error {
			return evolver(event.(E), k)
		}

		eType := reflect.TypeOf(proto)

		k.eventRoutes[eType] = append(k.eventRoutes[eType], handler)
	}
}
