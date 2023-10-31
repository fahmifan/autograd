package core

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Ctx struct {
	GormDB *gorm.DB
	JWTKey string
	MediaConfig
}

type MediaConfig struct {
	MediaServeBaseURL string
	RootFolder        string
	ObjectStorer      ObjectStorer
}

type ObjectStorer interface {
	Store(ctx context.Context, dst string, r io.Reader) error
}

func IsDBNotFoundErr(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func Transaction(ctx *Ctx, fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return ctx.GormDB.Transaction(fn)
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

func NewEntityMeta(now time.Time) TimestampMetadata {
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

func (p PaginationRequest) Offset() int32 {
	if p.Page <= 1 {
		return 0
	}

	return (p.Page - 1) * p.Limit
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
	return p.Total/p.Limit + 1
}

var (
	ProtoEmptyResponse = &connect.Response[autogradv1.Empty]{Msg: &autogradv1.Empty{}}
	ErrInternalServer  = connect.NewError(connect.CodeInternal, errors.New("internal server error"))
)
