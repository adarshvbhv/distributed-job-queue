package job

import "context"

type Handler interface {
	Execute(context context.Context, payload []byte) error
}


