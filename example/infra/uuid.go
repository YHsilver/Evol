package infra

import (
	"errors"
	"github.com/bwmarrin/snowflake"
)

func NewUUID() (int64, error) {
	n, err := snowflake.NewNode(1)
	if err != nil {
		return 0, errors.New("generate uuid failed")
	}

	return int64(n.Generate()), nil
}
