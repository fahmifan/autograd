package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/fahmifan/autograd/pkg/dbconn"
	"github.com/fahmifan/autograd/pkg/jobqueue"
	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/fahmifan/autograd/pkg/xsqlc"
	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
	"gorm.io/gorm"
)

var _mapHandlers = map[jobqueue.JobType]jobqueue.JobHandler{}

type HandlerFunc func(ctx context.Context, tx *gorm.DB, item jobqueue.OutboxItem) error

type EnqueueRequest struct {
	Payload       any
	JobType       jobqueue.JobType
	IdempotentKey jobqueue.IdempotentKey
}

type OutboxService struct {
	db        *gorm.DB
	sqlDB     *sql.DB
	debug     bool
	queuePool *queue.Queue

	stopChan chan bool
}

func NewOutboxService(db *gorm.DB, sqlDB *sql.DB, debug bool) *OutboxService {
	return &OutboxService{
		db:       db,
		debug:    debug,
		sqlDB:    sqlDB,
		stopChan: make(chan bool),
	}
}

func ValidJob(job jobqueue.JobType) bool {
	_, ok := _mapHandlers[job]
	return ok
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

	dbtx, ok := dbconn.DBTxFromGorm(tx)
	if !ok {
		return jobqueue.OutboxItem{}, logs.ErrWrapCtx(ctx, errors.New("transaction is invalid"), "OutboxService: Enqueue", "get dbtx")
	}

	if req.IdempotentKey != "" {
		oldItem, err := reader.FindPendingByKey(ctx, tx, string(req.IdempotentKey))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return jobqueue.OutboxItem{}, logs.ErrWrapCtx(ctx, err, "OutboxService: Enqueue", "find item")
		}

		hasPendingItem := oldItem.ID.String() != jobqueue.EmptyIDStr
		if hasPendingItem {
			return oldItem, nil
		}
	}

	item, err := jobqueue.NewOutboxItem(jobqueue.NewID(), req.JobType, req.IdempotentKey, payload)
	if err != nil {
		return jobqueue.OutboxItem{}, logs.ErrWrapCtx(ctx, err, "OutboxService: Enqueue", "new item")
	}

	err = writer.CreateV2(ctx, dbtx, &item)
	if err != nil {
		return jobqueue.OutboxItem{}, logs.ErrWrapCtx(ctx, err, "OutboxService: Enqueue", "save item to db")
	}

	return item, err
}

type QueueJob struct {
	Name    string
	Message string
}

func (job *QueueJob) Bytes() []byte {
	b, err := json.Marshal(job)
	if err != nil {
		panic(err)
	}
	return b
}

// Run will run blocking the OutboxService
func (svc *OutboxService) Run() error {
	const maxFetch = 100

	svc.queuePool = queue.NewPool(runtime.NumCPU(), queue.WithFn(func(ctx context.Context, m core.QueuedMessage) error {
		job, ok := m.(*QueueJob)
		if !ok {
			if err := json.Unmarshal(m.Bytes(), &job); err != nil {
				return err
			}
		}

		jobType := jobqueue.JobType(job.Name)
		for _, handler := range _mapHandlers {
			if handler.JobType() != jobType {
				continue
			}

			item := jobqueue.OutboxItem{}
			err := json.Unmarshal([]byte(job.Message), &item)
			if err != nil {
				logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "unmarshal item")
				continue
			}

			handle(svc.db, svc.sqlDB, svc.debug, handler)(ctx, svc.db, item)

			return nil
		}

		return nil
	}))
	defer func() {
		logs.Info("OutboxService: Run", "releasing jobqueue outbox queue pool")
		svc.queuePool.Release()
		logs.Info("OutboxService: Run", "done releasing jobqueue outbox queue pool")
	}()

	func() {
		for {
			// it's just an outbox should be fast
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			select {
			case <-svc.stopChan:
				logs.Info("OutboxService: Run", "stopping jobqueue outbox runner")
				cancel()
				return
			default:
				err := svc.run(ctx, maxFetch)
				if err != nil {
					logs.ErrCtx(ctx, err, "Run: OutboxService", "run")
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	logs.Info("OutboxService: Run", "done stopping jobqueue outbox runner")
	logs.Info("OutboxService: Run", "done stopping jobqueue")
	return nil
}

func (svc *OutboxService) Stop() {
	logs.Info("OutboxService: Run", "stopping jobqueue outbox")
	svc.stopChan <- true
}

func (svc *OutboxService) run(ctx context.Context, limit int) error {
	if svc.debug {
		logs.InfoCtx(ctx, "OutboxService: run", "start")
	}

	reader := OutboxItemReader{}

	ids, err := reader.FindAllPendingIDs(ctx, svc.sqlDB, limit)
	if err != nil {
		return logs.ErrWrapCtx(context.Background(), err, "Run: OutboxService", "find items")
	}

	if len(ids) == 0 {
		if svc.debug {
			logs.InfoCtx(ctx, "OutboxService: FindAllPendingIDs: empty", fmt.Sprint(ids))
		}
		return nil
	}

	if svc.debug {
		logs.InfoCtx(ctx, "OutboxService: FindAllPendingIDs", fmt.Sprint(ids))
	}

	delay := 2 * time.Millisecond
	idChan := make(chan jobqueue.ID, limit)
	wg := sync.WaitGroup{}
	// we need minimum worker of 2
	nworker := runtime.NumCPU()
	if nworker <= 0 {
		nworker = 2
	}

	workerFn := func(wg *sync.WaitGroup, workedID int, idChan <-chan jobqueue.ID) {
		defer wg.Done()

		if svc.debug {
			logs.InfoCtx(ctx, "OutboxService: run", "worker", fmt.Sprint(workedID))
		}

		for id := range idChan {
			logs.Info("OutboxService: run", "worker received id", id.String())
			err := svc.sendItem(ctx, workedID, id)
			if err != nil {
				logs.ErrCtx(ctx, err, "Run: OutboxService", "transaction", "itemID", id.String())
			}
			time.Sleep(delay)
		}
	}

	for i := range nworker {
		wg.Add(1)
		go workerFn(&wg, i, idChan)
	}

	for _, id := range ids {
		idChan <- id
	}

	close(idChan)
	wg.Wait()

	return nil
}

func (svc *OutboxService) sendItem(ctx context.Context, workerID int, id jobqueue.ID) error {
	outboxItemReader := OutboxItemReader{}
	writer := OutboxItemWriter{}

	err := dbconn.SqlcTransaction(ctx, svc.sqlDB, func(tx xsqlc.DBTX) error {
		item, err := outboxItemReader.FindByID(ctx, tx, id)
		if err != nil {
			return logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "find item")
		}

		item, err = item.MoveTo(jobqueue.StatusSent)
		if err != nil {
			return logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "move item")
		}

		err = writer.Update(ctx, tx, &item)
		if err != nil {
			return logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "update item")
		}

		itemBuf, err := json.Marshal(item)
		if err != nil {
			return logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "marshal item")
		}

		err = svc.queuePool.Queue(&QueueJob{
			Name:    string(item.JobType),
			Message: string(itemBuf),
		})
		if err != nil {
			logs.ErrWrapCtx(ctx, err, "Run: OutboxService", "queue job", string(item.JobType), "item", id.String())
			return err
		}

		if svc.debug {
			logs.InfoCtx(ctx, "OutboxService: run", "success queue job", string(item.JobType))
		}

		return nil
	})

	if err != nil {
		logs.ErrCtx(ctx, err, "Run: OutboxService", "transaction", "itemID", id.String())
		return err
	}

	if svc.debug {
		logs.InfoCtx(ctx, "OutboxService: run", "workerID", fmt.Sprint(workerID), "send item", id.String())
	}

	return nil
}

// RegisterHandlers register all job queue handler.
// This method is not thread safe, should be called only inside one goroutine.
func RegisterHandlers(db *gorm.DB, sqlDB *sql.DB, debug bool, handlers []jobqueue.JobHandler) {
	for _, handler := range handlers {
		_mapHandlers[handler.JobType()] = handler
	}
}

func handle(db *gorm.DB, _ *sql.DB, debug bool, handler jobqueue.JobHandler) HandlerFunc {
	return func(ctx context.Context, tx *gorm.DB, item jobqueue.OutboxItem) error {
		writer := OutboxItemWriter{}

		if debug {
			logs.InfoCtx(ctx, "outbox: handle", "item", string(item.JobType), "id", item.ID.String())
		}

		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
			dbtx, ok := dbconn.DBTxFromGorm(tx)
			if !ok {
				return logs.ErrWrapCtx(ctx, errors.New("transaction is invalid"), "outbox: handle: Get dbtx")
			}

			item, err = item.MoveTo(jobqueue.StatusPicked)
			if err != nil {
				return logs.ErrWrapCtx(ctx, err, "outbox: handle: MoveTo Picked")
			}

			err = writer.Update(ctx, dbtx, &item)
			if err != nil {
				return logs.ErrWrapCtx(ctx, err, "outbox: handle: Update Picked")
			}

			handleErr := handler.Handle(ctx, tx, item.Payload)
			if handleErr != nil {
				if item, err = item.MoveTo(jobqueue.StatusFailed); err != nil {
					return logs.ErrWrapCtx(ctx, fmt.Errorf("%w: %w", handleErr, err), "outbox: handle: MoveTo Failed")
				}

				if err = writer.Update(ctx, dbtx, &item); err != nil {
					return logs.ErrWrapCtx(ctx, err, "outbox: handle: Update Failed")
				}

				return logs.ErrWrapCtx(ctx, handleErr, "outbox: handle: Handle")
			}

			item, err = item.MoveTo(jobqueue.StatusSuccess)
			if err != nil {
				return logs.ErrWrapCtx(ctx, err, "outbox: handle: MoveTo Success")
			}

			if err = writer.Update(ctx, dbtx, &item); err != nil {
				return logs.ErrWrapCtx(ctx, err, "outbox: handle: Update Success")
			}

			return nil
		})
		if err != nil {
			return logs.ErrWrapCtx(ctx, err, "outbox: handle: Transaction")
		}

		if debug {
			logs.InfoCtx(ctx, "outbox: handle", "item", string(item.JobType), "id", item.ID.String(), "status", string(item.Status))
		}

		return nil
	}
}
