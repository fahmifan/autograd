package utils

import (
	"net/url"
	"sync"

	"github.com/miun173/autograd/model"
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

// GenerateCursorRequest generate a request for cursor
func GenerateCursorRequest(query url.Values) *model.CursorRequest {
	limit, page, sort := int64(10), int64(1), "created_at desc"

	for key, values := range query {
		value := values[0]

		switch key {
		case "limit":
			limit = StringToInt64(value)
			break
		case "page":
			page = StringToInt64(value)
			break
		case "sort":
			sort = value
			break
		}
	}

	return &model.CursorRequest{Limit: limit, Page: page, Sort: sort}
}
