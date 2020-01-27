package ui

type Events struct {
	events map[string][]func(string)
}

func (e Events) On(name string, f func(string)) {
	e.events[name] = append(e.events[name], f)
}

func (e Events) Emit(name string, value string) {
	es, ok := e.events[name]
	if ok {
		for _, f := range es {
			f(value)
		}
	}
}

func NewEvents() *Events {
	return &Events{
		events: make(map[string][]func(string)),
	}
}
