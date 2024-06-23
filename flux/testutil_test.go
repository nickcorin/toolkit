package flux_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nickcorin/toolkit/flux"
)

type testEvent struct {
	id        string
	topic     flux.EventTopic
	sequence  uint
	key       string
	timestamp time.Time
}

func (event *testEvent) ID() string             { return event.id }
func (event *testEvent) Topic() flux.EventTopic { return event.topic }
func (event *testEvent) Sequence() uint         { return event.sequence }
func (event *testEvent) Key() string            { return event.key }
func (event *testEvent) Timestamp() time.Time   { return event.timestamp }

type eventGenerator struct {
	t   *testing.T
	mu  sync.Mutex
	seq uint

	events []flux.Event
}

func NewEventGenerator(t *testing.T) *eventGenerator {
	t.Helper()

	return &eventGenerator{t: t}
}

func (g *eventGenerator) generateEvent(topic, key string) *testEvent {
	g.t.Helper()

	g.mu.Lock()
	defer g.mu.Unlock()

	g.seq++
	return &testEvent{
		id:        uuid.NewString(),
		topic:     flux.EventTopic(topic),
		sequence:  g.seq,
		key:       key,
		timestamp: time.Now().UTC(),
	}
}

func (g *eventGenerator) generateRandomEvent() *testEvent {
	g.t.Helper()

	return g.generateEvent(uuid.NewString(), uuid.NewString())
}

func (g *eventGenerator) Head(ctx context.Context) (flux.Event, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.events) == 0 {
		return nil, flux.ErrNoMoreEvents
	}

	return g.events[0], nil
}

func (g *eventGenerator) NextEvents(ctx context.Context, from, batchSize uint, streamLag time.Duration) ([]flux.Event, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i := 0; i < int(batchSize); i++ {
		g.events = append(g.events, g.generateRandomEvent())
	}

	return g.events[from : from+batchSize], nil
}
