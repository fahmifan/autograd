package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fahmifan/autograd/pkg/jobqueue"
	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/gookit/event"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type EnqueueRequest struct {
	Payload       any
	JobType       jobqueue.JobType
	IdempotentKey jobqueue.IdempotentKey
}

type OutboxService struct {
	db    *gorm.DB
	debug bool
}

func NewOutboxService(db *gorm.DB, debug bool) *OutboxService {
	return &OutboxService{
		db:    db,
		debug: debug,
	}
}

var _mapValidJob = map[jobqueue.JobType]bool{}

func registerValidJob(job jobqueue.JobType) {
	_mapValidJob[job] = true
}

func ValidJob(job jobqueue.JobType) bool {
	return _mapValidJob[job]
}

func (svc *OutboxService) Enqueue(ctx context.Context, tx *gorm.DB, req EnqueueRequest) (jobqueue.OutboxItem, error) {
	if !ValidJob(jobqueue.JobType(req.JobType)) {
		return jobqueue.OutboxItem{}, logs.ErrWrapCtx(ctx, errors.New("invalid destination"), "OutboxService: Enqueue", "valid job")
	}

	payload, err := jobqueue.MarshalPayload(req.Payload)
	if err != nil {
		return jobqueue.OutboxItem{}, logs.ErrWrapCtx(ctx, err, "OutboxService: Enqueue", "marshal body")
	}

	reader := OutboxItemReader{}
	writer := OutboxItemWriter{}

	oldItem, err := reader.FindPendingByKey(ctx, tx, string(req.IdempotentKey))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return jobqueue.OutboxItem{}, logs.ErrWrapCtx(ctx, err, "OutboxService: Enqueue", "find item")
	}

	hasPendingItem := oldItem.ID.String() != jobqueue.EmptyIDStr
	if hasPendingItem {
		return oldItem, nil
	}

	item, err := jobqueue.NewOutboxItem(jobqueue.NewID(), req.JobType, req.IdempotentKey, payload)
	if err != nil {
		return jobqueue.OutboxItem{}, logs.ErrWrapCtx(ctx, err, "OutboxService: Enqueue", "new item")
	}

	err = writer.Create(ctx, tx, item)
	if err != nil {
		return jobqueue.OutboxItem{}, logs.ErrWrapCtx(ctx, err, "OutboxService: Enqueue", "save item to db")
	}

	return item, err
}

// Run will run blocking the OutboxService
func (svc *OutboxService) Run(ctx context.Context) error {
	const maxFetch = 20

	for {
		time.Sleep(5 * time.Second)

		err := svc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return svc.run(ctx, tx, maxFetch)
		})
		if err != nil {
			logs.ErrCtx(ctx, err, "Run: OutboxService", "run")
		}
	}
}

func (svc *OutboxService) run(ctx context.Context, tx *gorm.DB, limit int) error {
	if svc.debug {
		logs.InfoCtx(ctx, "OutboxService: run", "start")
	}

	reader := OutboxItemReader{}

	items, err := reader.FindAllPending(ctx, tx, limit)
	if err != nil {
		return logs.ErrWrapCtx(context.Background(), err, "Run: OutboxService", "find items")
	}

	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		err, event := event.Fire(string(item.JobType), map[string]any{
			"item": item,
		})
		if err != nil {
			return logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "fire event", event.Name())
		}
		if svc.debug {
			logs.InfoCtx(ctx, "OutboxService: run", "fire items", event.Name())
		}
	}

	itemIDs := lo.Map(items, func(item jobqueue.OutboxItem, _ int) jobqueue.ID {
		return item.ID
	})

	writer := OutboxItemWriter{}

	err = writer.UpdateAllStatus(ctx, tx, items, jobqueue.StatusSent)
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "update items")
	}
	if svc.debug {
		logs.InfoCtx(ctx, "OutboxService: run", "update status: count:", fmt.Sprint(len(itemIDs)))
	}

	return nil
}

// RegisterHandlers register all job queue handler.
// This method is not thread safe, should be called only inside one goroutine.
func RegisterHandlers(db *gorm.DB, handlers []jobqueue.JobHandler) {
	for _, handler := range handlers {
		registerValidJob(handler.JobType())
		event.On(string(handler.JobType()), handle(db, handler))
	}
}

func handle(db *gorm.DB, handler jobqueue.JobHandler) event.ListenerFunc {
	return func(e event.Event) error {
		ctx := context.Background()
		writer := OutboxItemWriter{}

		item, ok := e.Get("item").(jobqueue.OutboxItem)
		if !ok {
			return logs.ErrWrapCtx(ctx, fmt.Errorf("invalid payload"), "outbox: handle: Get payload")
		}

		logs.InfoCtx(ctx, "outbox: handle", "item", string(item.JobType), "id", item.ID.String())

		items := []jobqueue.OutboxItem{item}
		err := writer.UpdateAllStatus(ctx, db, items, jobqueue.StatusPicked)
		if err != nil {
			return logs.ErrWrapCtx(ctx, err, "outbox: handle: Picked")
		}

		err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			payload := jobqueue.Payload(item.Payload)

			err = handler.Handle(ctx, tx, payload)
			if err != nil {
				return logs.ErrWrapCtx(ctx, err, "outbox: handle: Handle")
			}

			return nil
		})
		if err != nil {
			return logs.ErrWrapCtx(ctx, err, "outbox: handle: Transaction")
		}

		err = writer.UpdateAllStatus(ctx, db, items, jobqueue.StatusSuccess)
		if err != nil {
			// set status to failed
			if err = writer.UpdateAllStatus(ctx, db, items, jobqueue.StatusFailed); err != nil {
				return logs.ErrWrapCtx(ctx, err, "outbox: handle: set status Failed")
			}
		}

		return nil
	}
}
