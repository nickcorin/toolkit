package flux_test

import (
	"context"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/nickcorin/toolkit/flux"
	"github.com/nickcorin/toolkit/flux/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRelay_Shutdown(t *testing.T) {
	ctrl := gomock.NewController(t)

	generator := NewEventGenerator(t)

	dispatcher := mocks.NewMockDispatcher(ctrl)
	dispatcher.
		EXPECT().
		Dispatch(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	config := flux.DefaultRelayConfig
	req := flux.StreamRequest{}

	relay := flux.NewRelay(
		dispatcher,
		generator,
		flux.WithBufferSize(config.BufferSize),
		flux.WithBackOff(&backoff.ZeroBackOff{}),
	)

	go func() {
		err := relay.Start(context.Background(), req)
		require.NoError(t, err)
	}()

	relay.Shutdown()
}

func TestRelay_StreamEvents(t *testing.T) {
	t.Run("empty queue", func(t *testing.T) {
		t.Run("with default configs", func(t *testing.T) {
			config := flux.DefaultRelayConfig
			req := flux.StreamRequest{}

			testRelay_StreamEvents(t, 0, config, req)
		})
	})

	t.Run("non-empty queue", func(t *testing.T) {
		const eventCount = 50
		t.Run("with default configs", func(t *testing.T) {
			config := flux.DefaultRelayConfig
			req := flux.StreamRequest{}

			testRelay_StreamEvents(t, eventCount, config, req)
		})

		t.Run("with increased buffer size", func(t *testing.T) {
			config := flux.RelayConfig{BufferSize: 8}
			req := flux.StreamRequest{}

			testRelay_StreamEvents(t, eventCount, config, req)
		})

		t.Run("with start sequence", func(t *testing.T) {
			config := flux.DefaultRelayConfig
			req := flux.StreamRequest{StartSequence: 3}

			testRelay_StreamEvents(t, eventCount, config, req)
		})

		t.Run("with stream lag", func(t *testing.T) {
			config := flux.DefaultRelayConfig
			req := flux.StreamRequest{StreamLag: 5 * time.Minute}

			testRelay_StreamEvents(t, eventCount, config, req)
		})
	})
}

func setupEventReader(
	t *testing.T,
	ctrl *gomock.Controller,
	events []flux.Event,
	config flux.RelayConfig,
	req flux.StreamRequest,
) flux.EventReader {
	t.Helper()

	eventReader := mocks.NewMockEventReader(ctrl)

	// Filter out stream lagged events.
	i := 0
	for i < len(events) {
		if events[i].Timestamp().Add(req.StreamLag).After(time.Now()) {
			events = append(events[:i], events[i+1:]...)
			continue
		}

		i++
	}

	eventCount := uint(len(events))

	for cursor := req.StartSequence; cursor < eventCount; cursor += config.BufferSize {
		eventReader.
			EXPECT().
			NextEvents(gomock.Any(), cursor, config.BufferSize, req.StreamLag).
			Return(events[cursor:min(cursor+config.BufferSize, eventCount)], nil)
	}

	// Return no more events after the last batch.
	eventReader.
		EXPECT().
		NextEvents(gomock.Any(), uint(eventCount), config.BufferSize, req.StreamLag).
		Return(nil, flux.ErrNoMoreEvents).
		AnyTimes()

	return eventReader
}

func setupDispatcher(t *testing.T, ctrl *gomock.Controller, events []flux.Event, req flux.StreamRequest) flux.Dispatcher {
	t.Helper()

	dispatcher := mocks.NewMockDispatcher(ctrl)
	eventCount := uint(len(events))

	for i := uint(0); i < eventCount; i++ {
		if req.StartSequence >= events[i].Sequence() {
			continue
		}

		if time.Now().Add(-req.StreamLag).Before(events[i].Timestamp()) {
			continue
		}

		dispatcher.EXPECT().Dispatch(gomock.Any(), events[i]).Return(nil)
	}

	return dispatcher
}

func setupRelay(t *testing.T, ctrl *gomock.Controller, events []flux.Event, config flux.RelayConfig, req flux.StreamRequest) *flux.Relay {
	t.Helper()

	dispatcher := setupDispatcher(t, ctrl, events, req)
	eventReader := setupEventReader(t, ctrl, events, config, req)

	opts := []flux.RelayOption{
		flux.WithBufferSize(config.BufferSize),
		flux.WithBackOff(&backoff.StopBackOff{}),
	}

	relay := flux.NewRelay(dispatcher, eventReader, opts...)

	return relay
}

func testRelay_StreamEvents(t *testing.T, eventCount uint, relayOptions flux.RelayConfig, streamConfig flux.StreamRequest) {
	t.Helper()

	ctrl := gomock.NewController(t)

	g := NewEventGenerator(t)

	events := make([]flux.Event, eventCount)

	for i := uint(0); i < eventCount; i++ {
		e := g.generateRandomEvent()
		e.timestamp = e.timestamp.Add(-time.Duration(eventCount-1) * time.Minute)
		events[i] = e
	}

	relay := setupRelay(t, ctrl, events, relayOptions, streamConfig)

	err := relay.Start(context.Background(), streamConfig)
	require.ErrorIs(t, err, flux.ErrNoMoreEvents)

	relay.Shutdown()
}
