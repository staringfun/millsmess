package socket

import (
	"context"
	"encoding/json"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/types"
)

type Emitter interface {
	Emit(writer base.ContextWriter, event types.SocketMessageTypeEvent, data types.SocketData, ctx context.Context) error
}

type DefaultEmitter struct {
}

const EventDataArraySize = 2

func (e *DefaultEmitter) Emit(writer base.ContextWriter, event types.SocketMessageTypeEvent, data types.SocketData, ctx context.Context) error {
	bytes, err := json.Marshal([EventDataArraySize]any{event, data})
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes, ctx)
	return err
}
