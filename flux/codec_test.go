package flux_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nickcorin/toolkit/flux"
	"github.com/stretchr/testify/require"
)

func TestCodec(t *testing.T) {
	generator := NewEventGenerator(t)

	event := generator.generateEvent("test-topic", uuid.NewString())
	json := []byte("{" +
		"\"id\":\"" + event.ID() + "\"," +
		"\"topic\":\"" + string(event.Topic()) + "\"," +
		"\"sequence\":" + fmt.Sprintf("%d", event.Sequence()) + "," +
		"\"key\":\"" + event.Key() + "\"," +
		"\"timestamp\":\"" + event.Timestamp().Format(time.RFC3339Nano) + "\"" +
		"}")

	t.Run("json codec", func(t *testing.T) {
		codec := flux.NewJSONCodec()
		t.Run("encode", func(t *testing.T) {
			data, err := codec.Encode(event)
			require.NoError(t, err)
			require.NotEmpty(t, data)
			require.Equal(t, json, data)
		})

		t.Run("decode", func(t *testing.T) {
			e, err := codec.Decode(json)
			require.NoError(t, err)
			require.NotEmpty(t, e)

			require.EqualValues(t, event.ID(), e.ID())
			require.EqualValues(t, event.Topic(), e.Topic())
			require.EqualValues(t, event.Sequence(), e.Sequence())
			require.EqualValues(t, event.Key(), e.Key())

			// We use a different comparison for time.Time types since parsing through JSON loses the monotonic clock
			// field.
			require.True(t, event.Timestamp().Equal(e.Timestamp()))
		})
	})

	t.Run("protobuf codec", func(t *testing.T) {
		codec := flux.NewProtobufCodec()
		t.Run("encode", func(t *testing.T) {
			data, err := codec.Encode(event)
			require.NoError(t, err)
			require.NotEmpty(t, data)
		})

		t.Run("decode", func(t *testing.T) {
		})
	})
}
