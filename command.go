package kernel

type Command interface {
	Stream() string
}

type Dispatcher interface {
	Dispatch(cmd Command) error
}
