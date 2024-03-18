package outbox

import (
	"context"

	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/fahmifan/autograd/pkg/jobqueue"
	"github.com/fahmifan/autograd/pkg/xsqlc"
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
		return outboxItemFromModel(item)
	})

	return items, err
}

func (r *OutboxItemReader) FindAllPendingIDs(ctx context.Context, tx xsqlc.DBTX, limit int) (ids []jobqueue.ID, err error) {
	idStrs, err := xsqlc.New(tx).FindAllOutboxItemIDsByStatus(ctx, xsqlc.FindAllOutboxItemIDsByStatusParams{
		Status:    string(jobqueue.StatusPending),
		SizeLimit: int32(limit),
	})

	ids = lo.Map(idStrs, func(id string, _ int) jobqueue.ID {
		return mustParseID(id)
	})

	return ids, err
}

func (r *OutboxItemReader) FindPendingByKey(ctx context.Context, tx *gorm.DB, key string) (item jobqueue.OutboxItem, err error) {
	var outboxItem dbmodel.OutboxItem

	err = tx.Model(jobqueue.OutboxItem{}).
		Where("idempotent_key = ? AND status = ?", key, jobqueue.StatusPending).
		Take(&outboxItem).
		Error
	if err != nil {
		return jobqueue.OutboxItem{}, err
	}

	return outboxItemFromModel(outboxItem), nil
}

func (r *OutboxItemReader) FindPendingByKeyV2(ctx context.Context, tx xsqlc.DBTX, key string) (item jobqueue.OutboxItem, err error) {
	outboxItem, err := xsqlc.New(tx).FindOutboxItemByByKey(ctx, xsqlc.FindOutboxItemByByKeyParams{
		IdempotentKey: key,
		Status:        string(jobqueue.StatusPending),
	})

	return outboxItemFromSQLCModel(outboxItem), err
}

func (r *OutboxItemReader) FindByID(ctx context.Context, tx xsqlc.DBTX, id jobqueue.ID) (item jobqueue.OutboxItem, err error) {
	outboxItem, err := xsqlc.New(tx).FindOutboxItemByID(ctx, id.String())
	if err != nil {
		return jobqueue.OutboxItem{}, err
	}
	return outboxItemFromSQLCModel(outboxItem), err
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

func (r *OutboxItemWriter) CreateV2(ctx context.Context, tx xsqlc.DBTX, item *jobqueue.OutboxItem) error {
	outboxItem := dbmodel.OutboxItem{
		ID:            ulids.ULID(item.ID),
		JobType:       string(item.JobType),
		IdempotentKey: string(item.IdempotentKey),
		Status:        string(item.Status),
		Payload:       string(item.Payload),
	}

	res, err := xsqlc.New(tx).CreateOutboxItem(ctx, xsqlc.CreateOutboxItemParams{
		ID:            outboxItem.ID.String(),
		IdempotentKey: outboxItem.IdempotentKey,
		Status:        outboxItem.Status,
		JobType:       outboxItem.JobType,
		Payload:       outboxItem.Payload,
	})

	res.Version = outboxItem.Version

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

func (r *OutboxItemWriter) Update(ctx context.Context, tx xsqlc.DBTX, item *jobqueue.OutboxItem) error {
	res, err := xsqlc.New(tx).UpdateOutboxItem(ctx, xsqlc.UpdateOutboxItemParams{
		ID:            item.ID.String(),
		Status:        string(item.Status),
		IdempotentKey: string(item.IdempotentKey),
		JobType:       string(item.JobType),
		Payload:       string(item.Payload),
		Version:       int32(item.Version),
	})

	item.Version = res.Version

	return err
}

func outboxItemFromModel(model dbmodel.OutboxItem) jobqueue.OutboxItem {
	return jobqueue.OutboxItem{
		ID:            jobqueue.ID(model.ID),
		JobType:       jobqueue.JobType(model.JobType),
		IdempotentKey: jobqueue.IdempotentKey(model.IdempotentKey),
		Status:        jobqueue.Status(model.Status),
		Payload:       jobqueue.Payload(model.Payload),
		Version:       model.Version,
	}
}

func outboxItemFromSQLCModel(model xsqlc.OutboxItem) jobqueue.OutboxItem {
	return jobqueue.OutboxItem{
		ID:            mustParseID(model.ID),
		JobType:       jobqueue.JobType(model.JobType),
		IdempotentKey: jobqueue.IdempotentKey(model.IdempotentKey),
		Status:        jobqueue.Status(model.Status),
		Payload:       jobqueue.Payload(model.Payload),
		Version:       model.Version,
	}
}

func mustParseID(id string) jobqueue.ID {
	return jobqueue.ID(MustParseULID(id))
}

func MustParseULID(id string) ulids.ULID {
	ulid, err := ulids.Parse(id)
	if err != nil {
		panic(err)
	}
	return ulid
}
