package ui

type Events struct {
	events map[string][]func()
}

func (e Events) On(name string, f func()) {
	e.events[name] = append(e.events[name], f)
}

func (e Events) Emit(name string) {
	es, ok := e.events[name]
	if ok {
		for _, f := range es {
			f()
		}
	}
}

func NewEvents() *Events {
	return &Events{
		events: make(map[string][]func()),
	}
}
