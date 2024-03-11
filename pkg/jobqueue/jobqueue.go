package jobqueue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fahmifan/ulids"
	"gorm.io/gorm"
)

type Payload []byte
type JobType string

type JobHandler interface {
	Handle(ctx context.Context, tx *gorm.DB, payload Payload) error
	JobType() JobType
}

func MarshalPayload(v any) (Payload, error) {
	return json.Marshal(v)
}

func UnmarshalPayload(payload Payload, v any) error {
	return json.Unmarshal(payload, v)
}

const EmptyIDStr = "00000000000000000000000000"

type ID ulids.ULID
type Status string
type IdempotentKey string

const (
	StatusPending Status = "pending"
	StatusSent    Status = "sent"
	StatusPicked  Status = "picked"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

func NewID() ID {
	return ID(ulids.New())
}

type OutboxItem struct {
	ID            ID
	IdempotentKey IdempotentKey
	Status        Status
	JobType       JobType
	Payload       Payload
}

func NewOutboxItem(id ID, jobType JobType, key IdempotentKey, body Payload) (OutboxItem, error) {
	if id.String() == "" {
		return OutboxItem{}, errors.New("invalid id")
	}

	if jobType == "" {
		return OutboxItem{}, errors.New("invalid destination")
	}

	if key == "" {
		key = IdempotentKey(id.String())
	}

	item := OutboxItem{
		ID:            id,
		JobType:       JobType(jobType),
		Payload:       body,
		Status:        StatusPending,
		IdempotentKey: key,
	}

	return item, nil
}
