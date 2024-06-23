package flux

import (
	"context"
	"errors"
	"time"
)

type EventTopic string

func (topic EventTopic) String() string {
	return string(topic)
}

// ErrEventNotFound is an error that is returned when an event cannot be found.
var ErrEventNotFound = errors.New("event not found")

// ErrNoMoreEvents is an error that is returned when there are no event with a higher sequence number than provided.
var ErrNoMoreEvents = errors.New("no more events")

type Event interface {
	// ID returns a unique identifier for the event.
	ID() string

	// Topic is a partitioning key that can be used to filter events.
	Topic() EventTopic

	// Sequence returns an unsigned integer that can be used to order events. Events with a lower sequence number should
	// be considered to have been emitted before those with a higher sequence number.
	Sequence() uint

	// Key is a reference to the entity that the event is related to.
	Key() string

	// Timestamp returns the time at which the event was emitted.
	Timestamp() time.Time
}

// Compile-time assertion that defaultEvent implements the Event interface.
var _ Event = (*defaultEvent)(nil)

type defaultEvent struct {
	id        string
	topic     EventTopic
	sequence  uint
	key       string
	timestamp time.Time
}

func (event *defaultEvent) ID() string           { return event.id }
func (event *defaultEvent) Topic() EventTopic    { return event.topic }
func (event *defaultEvent) Sequence() uint       { return event.sequence }
func (event *defaultEvent) Key() string          { return event.key }
func (event *defaultEvent) Timestamp() time.Time { return event.timestamp }

// EventStore is an interface that combines an EventReader and an EventWriter.
type EventStore interface {
	EventReader
	EventWriter
}

// EventReader allows read-only access to an event store.
type EventReader interface {
	// Head returns the most recent event in the event store.
	//
	// Must return ErrEventNotFound if the event store is empty.
	Head(ctx context.Context) (Event, error)

	// NextEvents returns the next batch of events in the event store.
	//
	// Must return ErrNoMoreEvents if there are no events in the event store.
	NextEvents(ctx context.Context, from, batchSize uint, streamLag time.Duration) ([]Event, error)
}

// EventWriter allows write-only access to an event store.
type EventWriter interface {
	// CreateEvent creates a new event in the event store.
	CreateEvent(ctx context.Context, topic, key string) (Event, error)
}

// EventFilter defines a type that determines whether an event should be included in the event stream.
//
// Apply should return true if the event should be included in the event stream.
type EventFilter interface {
	Apply(e Event) bool
}

// EventFilterFunc is a function type that implements the EventFilter interface.
type EventFilterFunc func(e Event) bool

// Filter returns true if the event should be included in the event stream.
func (f EventFilterFunc) Apply(e Event) bool {
	return f(e)
}

// MatchKey returns an EventFilter that filters events by their keys.
func MatchKey(key string) EventFilterFunc {
	return func(e Event) bool {
		return e.Key() == key
	}
}

// MatchTopics returns an EventFilter that filters events by their topics.
func MatchTopics(topics ...EventTopic) EventFilterFunc {
	return func(e Event) bool {
		for _, topic := range topics {
			if e.Topic() == topic {
				return true
			}
		}

		return false
	}
}
