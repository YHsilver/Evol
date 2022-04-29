package codec

import (
	"context"
	"encoding/json"
	"evol"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"time"
)

type JsonEventCodec struct {
}

func (j *JsonEventCodec) MarshalEvent(ctx context.Context, e evol.Event) ([]byte, error) {
	newEvent := &evt{
		Topic:         e.Topic(),
		Data:          e.Data(),
		AggregateType: e.AggregateType(),
		AggregateId:   e.AggregateIdentity(),
		Time:          e.Timestamp(),
	}
	return json.Marshal(newEvent)
}

func (j *JsonEventCodec) UnmarshalEvent(ctx context.Context, bytes []byte) (evol.Event, error) {
	var e evt
	if err := json.Unmarshal(bytes, &e); err != nil {
		return nil, fmt.Errorf("unmarshal e error: %w", err)
	}
	res := evol.NewEvent(
		e.Topic,
		e.Data,
		e.Time,
		evol.ForAggregate(e.AggregateType, e.AggregateId),
	)

	return res, nil
}

type evt struct {
	Topic         evol.Topic         `json:"topic"`
	Data          interface{}        `json:"data"`
	AggregateType evol.AggregateType `json:"aggregate_type"`
	AggregateId   string             `json:"aggregate_id"`
	Time          time.Time          `json:"time"`
}

// DecodeEventData translate map[string]interface{} to struct
// because json.Unmarshal will decode interface{} to map[string]interface{}, thus it's useful for decode codec data to certain struct
func DecodeEventData(rawData interface{}, output interface{}) error {
	return mapstructure.Decode(rawData, output)
}
