package outbox

// import (
// 	"context"
// 	"errors"
// 	"time"

// 	"github.com/fahmifan/autograd/pkg/logs"
// 	"github.com/fahmifan/ulids"
// 	"github.com/gookit/event"
// 	"github.com/samber/lo"
// 	"gorm.io/gorm"
// )

// type ID ulids.ULID
// type Destination string
// type Status string
// type IdempotentKey string

// const (
// 	StatusPending Status = "pending"
// 	StatusSent    Status = "sent"
// 	StatusPicked  Status = "picked"
// 	StatusFailed  Status = "failed"
// )

// func NewID() ID {
// 	return ID(ulids.New())
// }

// type Item struct {
// 	ID            ID
// 	IdempotentKey IdempotentKey
// 	Status        Status
// 	Destination   Destination
// 	Body          string
// }

// func NewItem(id ID, dest Destination, key IdempotentKey, body string) (Item, error) {
// 	if id.String() == "" {
// 		return Item{}, errors.New("invalid id")
// 	}

// 	if dest == "" {
// 		return Item{}, errors.New("invalid destination")
// 	}

// 	if key == "" {
// 		key = IdempotentKey(id.String())
// 	}

// 	item := Item{
// 		ID:            id,
// 		Destination:   dest,
// 		Body:          body,
// 		Status:        StatusPending,
// 		IdempotentKey: key,
// 	}

// 	return item, nil
// }

// type EnqueueRequest struct {
// 	Body          string
// 	Destination   Destination
// 	IdempotentKey IdempotentKey
// }

// type OutboxService struct {
// 	db *gorm.DB
// }

// func NewOutboxService(db *gorm.DB) OutboxService {
// 	return OutboxService{
// 		db: db,
// 	}
// }

// func (svc *OutboxService) Enqueue(ctx context.Context, tx *gorm.DB, req EnqueueRequest) (Item, error) {
// 	item, err := NewItem(NewID(), req.Destination, req.IdempotentKey, req.Body)
// 	if err != nil {
// 		return Item{}, logs.ErrWrapCtx(ctx, err, "OutboxService: Enqueue", "new item")
// 	}

// 	err = tx.WithContext(ctx).Create(&item).Error
// 	if err != nil {
// 		return Item{}, logs.ErrWrapCtx(ctx, err, "OutboxService: Enqueue", "save item to db")
// 	}

// 	return item, err
// }

// // Run will run blocking the OutboxService
// func (svc *OutboxService) Run(ctx context.Context) error {
// 	const maxFetch = 1

// 	for {
// 		time.Sleep(5 * time.Millisecond)

// 		err := svc.db.Transaction(func(tx *gorm.DB) error {
// 			return svc.run(ctx, tx, maxFetch)
// 		})
// 		if err != nil {
// 			logs.ErrCtx(ctx, err, "Run: OutboxService", "process")
// 		}
// 	}
// }

// func (svc *OutboxService) run(ctx context.Context, tx *gorm.DB, limit int) error {
// 	var items []Item
// 	err := tx.Model(Item{}).
// 		Where("status = ?", StatusPending).
// 		Order("id ASC").
// 		Limit(limit).
// 		Find(&items).
// 		Error
// 	if err != nil {
// 		return logs.ErrWrapCtx(context.Background(), err, "Run: OutboxService", "find items")
// 	}

// 	if len(items) == 0 {
// 		return nil
// 	}

// 	for _, item := range items {
// 		err, _ := event.Fire(string(item.Destination), map[string]any{
// 			"item": item,
// 		})
// 		if err != nil {
// 			return logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "fire event")
// 		}
// 	}

// 	itemIDs := lo.Map(items, func(item Item, _ int) ID {
// 		return item.ID
// 	})

// 	err = tx.Model(Item{}).
// 		Where("id IN ? AND status = ?", itemIDs, StatusPending).
// 		Update("status", StatusSent).
// 		Error
// 	if err != nil {
// 		return logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "update items")
// 	}

// 	return nil
// }
