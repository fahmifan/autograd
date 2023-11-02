package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"

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

// GenerateUniqueString ..
func GenerateUniqueString() string {
	h := md5.New()
	randomNumber := fmt.Sprint(rand.Intn(10))
	timestamp := fmt.Sprint(time.Now().Unix())

	h.Write([]byte(randomNumber + timestamp))

	return hex.EncodeToString(h.Sum(nil))
}
