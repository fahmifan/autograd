package core

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/fahmifan/autograd/pkg/jobqueue"
	"github.com/fahmifan/autograd/pkg/jobqueue/outbox"
	"github.com/fahmifan/autograd/pkg/mailer"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Ctx struct {
	MediaConfig
	Debug       bool
	JWTKey      string
	SenderEmail string
	AppLink     string
	LogoURL     string

	GormDB         *gorm.DB
	SqlDB          *sql.DB
	Mailer         mailer.Mailer
	OutboxEnqueuer OutboxEnqueuer
}

type OutboxEnqueuer interface {
	Enqueue(ctx context.Context, tx *gorm.DB, req outbox.EnqueueRequest) (item jobqueue.OutboxItem, err error)
}

type MediaConfig struct {
	MediaServeBaseURL string
	RootDir           string
	ObjectStorer      ObjectStorer
}

type ObjectStorer interface {
	Store(ctx context.Context, dst string, r io.Reader) error
	Seek(ctx context.Context, srcpath string) (io.ReadCloser, error)
}

func IsDBNotFoundErr(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func Transaction(ctx context.Context, coreCtx *Ctx, fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return coreCtx.GormDB.WithContext(ctx).Transaction(fn)
}

type TimestampMetadata struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt null.Time
}

func (meta TimestampMetadata) ModelMetadata() dbmodel.Metadata {
	return dbmodel.Metadata{
		CreatedAt: null.TimeFrom(meta.CreatedAt),
		UpdatedAt: null.TimeFrom(meta.UpdatedAt),
		DeletedAt: gorm.DeletedAt(meta.DeletedAt.NullTime),
	}
}

func (meta TimestampMetadata) CreatedAtRFC3339() string {
	return meta.CreatedAt.Format(time.RFC3339)
}

func (meta TimestampMetadata) UpdatedAtRFC3339() string {
	return meta.CreatedAt.Format(time.RFC3339)
}

func (meta TimestampMetadata) DeletedAtRFC3339() null.String {
	if !meta.DeletedAt.Valid {
		return null.String{}
	}
	return null.StringFrom(meta.DeletedAt.Time.Format(time.RFC3339))
}

func (meta TimestampMetadata) ProtoTimestampMetadata() *autogradv1.TimestampMetadata {
	return &autogradv1.TimestampMetadata{
		CreatedAt: meta.CreatedAtRFC3339(),
		UpdatedAt: meta.UpdatedAtRFC3339(),
	}
}

func NewTimestampMeta(now time.Time) TimestampMetadata {
	return TimestampMetadata{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewModelMetadata(meta TimestampMetadata) dbmodel.Metadata {
	return dbmodel.Metadata{
		CreatedAt: null.TimeFrom(meta.CreatedAt),
		UpdatedAt: null.TimeFrom(meta.UpdatedAt),
		DeletedAt: gorm.DeletedAt(meta.DeletedAt.NullTime),
	}
}

func TimestampMetaFromModel(meta dbmodel.Metadata) TimestampMetadata {
	return TimestampMetadata{
		CreatedAt: meta.CreatedAt.Time,
		UpdatedAt: meta.UpdatedAt.Time,
		DeletedAt: null.Time{NullTime: sql.NullTime(meta.DeletedAt)},
	}
}

type PaginationRequest struct {
	Page  int32
	Limit int32
}

func PaginationRequestFromProto(p *autogradv1.PaginationRequest) PaginationRequest {
	return PaginationRequest{
		Page:  p.GetPage(),
		Limit: p.GetLimit(),
	}
}

func (p PaginationRequest) Offset() int32 {
	if p.Page <= 1 {
		return 0
	}

	return (p.Page - 1) * p.Limit
}

func (p PaginationRequest) PaginateScope(tx *gorm.DB) *gorm.DB {
	return tx.Limit(int(p.Limit)).Offset(int(p.Offset()))
}

type Pagination struct {
	Page  int32
	Limit int32
	Total int32
}

func (p Pagination) ProtoPagination() *autogradv1.PaginationMetadata {
	return &autogradv1.PaginationMetadata{
		Total:     p.Total,
		Page:      p.Page,
		Limit:     p.Limit,
		TotalPage: p.TotalPage(),
	}
}

func (p Pagination) Offset() int32 {
	if p.Page <= 1 {
		return 0
	}

	return (p.Page - 1) * p.Limit
}

func (p Pagination) TotalPage() int32 {
	if p.Total < p.Limit {
		return 1
	}

	if p.Total%p.Limit == 0 {
		return p.Total / p.Limit
	}

	return p.Total/p.Limit + 1
}

var (
	ProtoEmptyResponse  = &connect.Response[autogradv1.Empty]{Msg: &autogradv1.Empty{}}
	ErrInternalServer   = connect.NewError(connect.CodeInternal, errors.New("internal server error"))
	ErrUnauthenticated  = connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	ErrPermissionDenied = connect.NewError(connect.CodePermissionDenied, errors.New("unauthorized"))
)
