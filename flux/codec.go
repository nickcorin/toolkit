package flux

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nickcorin/toolkit/flux/fluxpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Codec is a message encode and decoder.
type Codec interface {
	Encode(e Event) ([]byte, error)
	Decode(data []byte) (Event, error)
}

// NewJSONCodec returns a new JSON codec.
func NewJSONCodec() Codec {
	return &JSONCodec{}
}

type JSONCodec struct{}

type jsonEvent struct {
	ID        string     `json:"id"`
	Topic     EventTopic `json:"topic"`
	Sequence  uint       `json:"sequence"`
	Key       string     `json:"key"`
	Timestamp time.Time  `json:"timestamp"`
}

func (codec *JSONCodec) Encode(e Event) ([]byte, error) {
	jsonEvent := jsonEvent{
		ID:        e.ID(),
		Topic:     e.Topic(),
		Sequence:  e.Sequence(),
		Key:       e.Key(),
		Timestamp: e.Timestamp(),
	}

	return json.Marshal(&jsonEvent)
}

func (codec *JSONCodec) Decode(data []byte) (Event, error) {
	var jsonEvent jsonEvent
	err := json.Unmarshal(data, &jsonEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	e := defaultEvent{
		id:        jsonEvent.ID,
		topic:     jsonEvent.Topic,
		sequence:  jsonEvent.Sequence,
		key:       jsonEvent.Key,
		timestamp: jsonEvent.Timestamp,
	}

	return &e, nil
}

// NewProtobufCodec returns a new Protobuf codec.
func NewProtobufCodec() Codec {
	return &ProtobufCodec{}
}

type ProtobufCodec struct{}

func (codec *ProtobufCodec) Encode(e Event) ([]byte, error) {
	return proto.Marshal(EventToProto(e))
}

func (codec *ProtobufCodec) Decode(data []byte) (Event, error) {
	var pb fluxpb.Event
	if err := proto.Unmarshal(data, &pb); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return EventFromProto(&pb), nil
}

func EventFromProto(pb *fluxpb.Event) Event {
	return &defaultEvent{
		id:        pb.Id,
		topic:     EventTopic(pb.Topic),
		sequence:  uint(pb.Sequence),
		key:       pb.Key,
		timestamp: pb.Timestamp.AsTime(),
	}
}

func EventToProto(e Event) *fluxpb.Event {
	return &fluxpb.Event{
		Id:        e.ID(),
		Topic:     string(e.Topic()),
		Sequence:  uint64(e.Sequence()),
		Key:       e.Key(),
		Timestamp: timestamppb.New(e.Timestamp()),
	}
}
