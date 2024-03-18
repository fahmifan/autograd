package jobqueue

import (
	"context"
	"encoding/json"
	"errors"
	"slices"

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
	Version       int32
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

func (item OutboxItem) IsEmpty() bool {
	return item.ID.String() == EmptyIDStr
}

type From = Status
type To = Status

var _statusFSM = map[From][]To{
	StatusPending: {StatusSent},
	StatusSent:    {StatusPicked},
	StatusPicked:  {StatusSuccess, StatusFailed},
	StatusFailed:  nil,
	StatusSuccess: nil,
}

func (item OutboxItem) MoveTo(nextStatus Status) (OutboxItem, error) {
	if !item.canTransitionTo(nextStatus) {
		return item, errors.New("invalid status transition")
	}
	item.Status = nextStatus
	return item, nil
}

func (item OutboxItem) canTransitionTo(nextStatus Status) bool {
	allowedStatues := _statusFSM[item.Status]
	if len(allowedStatues) == 0 {
		return false
	}

	return slices.Contains(allowedStatues, nextStatus)
}
