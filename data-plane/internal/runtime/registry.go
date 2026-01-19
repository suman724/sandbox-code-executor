package runtime

type Adapter interface {
	Run(code string) error
}

type Registry struct {
	Adapters map[string]Adapter
}

func NewRegistry() Registry {
	return Registry{Adapters: map[string]Adapter{}}
}

func (r Registry) Register(language string, adapter Adapter) {
	if r.Adapters == nil {
		r.Adapters = map[string]Adapter{}
	}
	r.Adapters[language] = adapter
}

func (r Registry) Adapter(language string) (Adapter, bool) {
	adapter, ok := r.Adapters[language]
	return adapter, ok
}

func (r Registry) Supports(language string) bool {
	_, ok := r.Adapters[language]
	return ok
}
