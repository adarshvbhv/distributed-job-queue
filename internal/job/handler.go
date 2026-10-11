package job

import "context"

type Handler interface {
	Execute(context context.Context, payload []byte) error
}


type Registry struct{
	handlers map[string]Handler
}

func NewRegistry () *Registry{
	return &Registry{
		handlers: make(map[string]Handler),
	}
}

func (r *Registry) Register(jobType string, handler Handler) {
	r.handlers[jobType] = handler
}

func (r *Registry) Get(jobType string) (Handler, bool){
	handler, ok:= r.handlers[jobType]

	return handler, ok

}