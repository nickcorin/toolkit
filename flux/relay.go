package flux

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
)

// Relay is a message broker that reads events from an event store and passes them to a dispatcher.
type Relay struct {
	dispatcher Dispatcher
	events     EventReader

	config *RelayConfig

	running  bool
	shutdown chan struct{}
}

type RelayConfig struct {
	// The backoff strategy to use when retrying operations.
	BackOff backoff.BackOff

	// The size of the buffer used to store events before they are dispatched.
	BufferSize uint
}

var DefaultRelayConfig = RelayConfig{
	BackOff:    backoff.NewConstantBackOff(5 * time.Second),
	BufferSize: 5,
}

type StreamRequest struct {
	// The sequence to start streaming events from. If zero, the relay will start streaming from the lowest sequence in
	// the event store.
	StartSequence uint

	// The filters to apply to the event stream. If empty, the relay will stream all events.
	Filters []EventFilter

	// If set, the relay will only stream events which occurred before NOW - StreamLag.
	StreamLag time.Duration
}

// NewRelay returns an instance of a Relay.
func NewRelay(dispatcher Dispatcher, events EventReader, opts ...RelayOption) *Relay {
	defaultConfig := DefaultRelayConfig // make a copy so we don't modify the default config.

	r := Relay{
		dispatcher: dispatcher,
		events:     events,
		config:     &defaultConfig,
		running:    false,
		shutdown:   make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt.Apply(r.config)
	}

	return &r
}

// RelayOption is an interface that allows for functional options to be applied to a RelayConfig.
type RelayOption interface {
	Apply(*RelayConfig)
}

// RelayOptionFunc is a function type that implements the RelayOption interface.
type RelayOptionFunc func(*RelayConfig)

// Apply applies the function to the relay.
func (f RelayOptionFunc) Apply(config *RelayConfig) {
	f(config)
}

// WithBufferSize sets the buffer size of the relay.
func WithBufferSize(bufferSize uint) RelayOption {
	return RelayOptionFunc(func(config *RelayConfig) {
		if bufferSize > 0 {
			config.BufferSize = bufferSize
		}
	})
}

// WithBackOff sets the backoff strategy of the relay.
func WithBackOff(backOff backoff.BackOff) RelayOption {
	return RelayOptionFunc(func(config *RelayConfig) {
		config.BackOff = backOff
	})
}

func (r *Relay) Start(ctx context.Context, req StreamRequest) error {
	if r.running {
		return fmt.Errorf("relay is already running")
	}

	r.running = true

	fn := func() error {
		return r.stream(ctx, req)
	}

	notify := func(err error, next time.Duration) {
		log.Printf("error: %v, next: %v", err, next)
	}

	return backoff.RetryNotify(fn, backoff.WithContext(r.config.BackOff, ctx), notify)
}

func (r *Relay) stream(ctx context.Context, req StreamRequest) error {
	if !r.running {
		return backoff.Permanent(fmt.Errorf("relay is not running"))
	}

	s := &stream{
		buffer:     make([]Event, 0),
		bufferSize: r.config.BufferSize,
		filters:    req.Filters,
		lag:        req.StreamLag,
		position:   req.StartSequence,
	}

	// Ensure all events are dispatched before returning.
	defer s.flush(ctx, r.dispatcher)

	// Start main loop.
	for {
		select {
		case <-r.shutdown:
			return nil
		default:
			if err := s.spool(ctx, r.events); err != nil {
				return fmt.Errorf("failed to spool events: %w", err)
			}

			if err := s.flush(ctx, r.dispatcher); err != nil {
				return fmt.Errorf("failed to flush events: %w", err)
			}

			// Reset the backoff strategy if we successfully processed an event batch.
			r.config.BackOff.Reset()
		}
	}
}

// Shutdown gracefully shuts down the relay, flushing any events that are still in the buffer.
func (r *Relay) Shutdown() {
	if !r.running {
		return
	}

	r.running = false
	close(r.shutdown)
}

type stream struct {
	bufferSize uint
	buffer     []Event
	filters    []EventFilter
	lag        time.Duration
	position   uint
}

func (s *stream) maybeDispatchEvent(ctx context.Context, e Event, d Dispatcher) error {
	// Detect gaps in the event stream. Note, this must be run before filtering.
	if e.Sequence() != s.position+1 {
		return fmt.Errorf("gap detected between event %d and %d", s.position, e.Sequence())
	}

	// Apply filters.
	shouldDispatch := true
	for _, filter := range s.filters {
		if !filter.Apply(e) {
			shouldDispatch = false
		}
	}

	if shouldDispatch {
		if err := d.Dispatch(ctx, e); err != nil {
			return fmt.Errorf("failed to dispatch event %s: %w", e.ID(), err)
		}
	}

	s.position = e.Sequence()

	return nil
}

func (s *stream) spool(ctx context.Context, events EventReader) error {
	el, err := events.NextEvents(ctx, s.position, s.bufferSize, s.lag)
	if err != nil {
		return fmt.Errorf("failed to fetch next events: %w", err)
	}

	s.buffer = append(s.buffer, el...)

	return nil
}

func (s *stream) flush(ctx context.Context, d Dispatcher) error {
	for len(s.buffer) > 0 {
		e := s.buffer[0]
		s.buffer = s.buffer[1:]

		if err := s.maybeDispatchEvent(ctx, e, d); err != nil {
			return err
		}
	}

	return nil
}
