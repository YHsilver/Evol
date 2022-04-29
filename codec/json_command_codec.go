package codec

import (
	"context"
	"encoding/json"
	"evol"
)

type JsonCmdCodec struct {
}

func (j *JsonCmdCodec) MarshalCommand(ctx context.Context, cmd evol.Command) ([]byte, error) {
	c := command{
		Name:                    cmd.Name(),
		TargetAggregateType:     cmd.TargetAggregateType(),
		TargetAggregateIdentity: cmd.TargetIdentity(),
	}
	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	c.Data = data

	return json.Marshal(c)
}

func (j *JsonCmdCodec) UnmarshalCommand(ctx context.Context, bytes []byte) (evol.Command, error) {
	var c command
	err := json.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	c2 := evol.NewCommand(c.Name)

	json.Unmarshal(c.Data, &c2)

	return c2, nil

}

type command struct {
	Name                    evol.CommandName   `json:"command_name,omitempty"`
	TargetAggregateType     evol.AggregateType `json:"target_aggregate_type,omitempty"`
	TargetAggregateIdentity string             `json:"target_aggregate_identity,omitempty"`
	Data                    json.RawMessage    `json:"data,omitempty"`
}
