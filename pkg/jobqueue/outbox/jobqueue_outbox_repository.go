package outbox

import (
	"context"

	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/fahmifan/autograd/pkg/jobqueue"
	"github.com/fahmifan/ulids"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type OutboxItemReader struct{}

func (r *OutboxItemReader) FindAllPending(ctx context.Context, tx *gorm.DB, limit int) (items []jobqueue.OutboxItem, err error) {
	var outboxItems []dbmodel.OutboxItem

	err = tx.Model(jobqueue.OutboxItem{}).
		Where("status = ?", jobqueue.StatusPending).
		Order("id ASC").
		Limit(limit).
		Find(&outboxItems).
		Error

	items = lo.Map(outboxItems, func(item dbmodel.OutboxItem, _ int) jobqueue.OutboxItem {
		return jobqueue.OutboxItem{
			ID:            jobqueue.ID(item.ID),
			JobType:       jobqueue.JobType(item.JobType),
			IdempotentKey: jobqueue.IdempotentKey(item.IdempotentKey),
			Status:        jobqueue.Status(item.Status),
			Payload:       jobqueue.Payload(item.Payload),
		}
	})

	return items, err
}

type OutboxItemWriter struct{}

func (r *OutboxItemWriter) Create(ctx context.Context, tx *gorm.DB, item jobqueue.OutboxItem) error {
	outboxItem := dbmodel.OutboxItem{
		ID:            ulids.ULID(item.ID),
		JobType:       string(item.JobType),
		IdempotentKey: string(item.IdempotentKey),
		Status:        string(item.Status),
		Payload:       string(item.Payload),
	}

	err := tx.Create(&outboxItem).Error
	return err
}

func (r *OutboxItemWriter) UpdateAllStatus(ctx context.Context, tx *gorm.DB, items []jobqueue.OutboxItem, status jobqueue.Status) error {
	var ids []string
	for _, item := range items {
		ids = append(ids, item.ID.String())
	}

	err := tx.Model(jobqueue.OutboxItem{}).
		Where("id IN ?", ids).
		Update("status", status).
		Error

	return err
}
