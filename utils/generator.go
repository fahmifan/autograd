package utils

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake
var generatorOnce sync.Once

func init() {
	generatorOnce.Do(func() {
		flake = sonyflake.NewSonyflake(sonyflake.Settings{})
	})
}

// GenerateID generate unique int64 id
func GenerateID() int64 {
	id, err := flake.NextID()
	if err != nil {
		logrus.Error(err)
	}

	return int64(id)
}
