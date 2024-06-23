package flux

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nickcorin/toolkit/flux/fluxpb"
)

// Dispatcher is a type that is responsible for dispatching events to various network components, or message queues.
type Dispatcher interface {
	Dispatch(ctx context.Context, e Event) error
}

// Compile-time assertion that GRPCDispatcher implements the Dispatcher interface.
var _ Dispatcher = (*GRPCDispatcher)(nil)

// NewGRPCDispatcher returns a new gRPC dispatcher.
func NewGRPCDispatcher(conn fluxpb.Flux_DispatchServer) *GRPCDispatcher {
	return &GRPCDispatcher{
		conn: conn,
	}
}

// GRPCDispatcher is a dispatcher that sends events to a gRPC server.
type GRPCDispatcher struct {
	conn fluxpb.Flux_DispatchServer
}

func (d *GRPCDispatcher) Dispatch(ctx context.Context, e Event) error {
	return d.conn.Send(EventToProto(e))
}

// Compile-time assertion that NatsDispatcher implements the Dispatcher interface.
var _ Dispatcher = (*NatsDispatcher)(nil)

// NewNatsDispatcher returns a new NATS dispatcher.
func NewNatsDispatcher(conn *nats.Conn, codec Codec) *NatsDispatcher {
	return &NatsDispatcher{
		codec: codec,
		conn:  conn,
	}
}

// NatsDispatcher is a dispatcher that sends events to a NATS message broker.
type NatsDispatcher struct {
	codec Codec
	conn  *nats.Conn
}

func (d *NatsDispatcher) Dispatch(ctx context.Context, e Event) error {
	msg, err := d.codec.Encode(e)
	if err != nil {
		return fmt.Errorf("failed to encode event: %w", err)
	}

	if err := d.conn.Publish(e.Topic().String(), msg); err != nil {
		return fmt.Errorf("failed to dispatch event: %w", err)
	}

	return nil
}
