package kernel

type State interface {
	On(Event) bool
}
