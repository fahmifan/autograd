package core

import (
	"database/sql"
	"errors"
	"time"

	"connectrpc.com/connect"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Ctx struct {
	GormDB *gorm.DB
}

func IsDBNotFoundErr(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func Transaction(ctx *Ctx, fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return ctx.GormDB.Transaction(fn)
}

type EntityMeta struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt null.Time
}

func (meta EntityMeta) CreatedAtRFC3339() string {
	return meta.CreatedAt.Format(time.RFC3339)
}

func (meta EntityMeta) UpdatedAtRFC3339() string {
	return meta.CreatedAt.Format(time.RFC3339)
}

func (meta EntityMeta) DeletedAtRFC3339() null.String {
	if !meta.DeletedAt.Valid {
		return null.String{}
	}
	return null.StringFrom(meta.DeletedAt.Time.Format(time.RFC3339))
}

func (meta EntityMeta) ProtoMeta() *autogradv1.Metadata {
	return &autogradv1.Metadata{
		CreatedAt: meta.CreatedAtRFC3339(),
		UpdatedAt: meta.UpdatedAtRFC3339(),
	}
}

func NewEntityMeta(now time.Time) EntityMeta {
	return EntityMeta{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

var EmptyResponse = &connect.Response[autogradv1.Empty]{
	Msg: &autogradv1.Empty{},
}
