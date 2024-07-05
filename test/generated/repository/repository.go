package repository

import (
	"context"

	helper "github.com/louvri/gosl/builder"
	"github.com/louvri/gosl/transformer"
)

type Repository interface {
	Get(ctx context.Context, param helper.QueryParams, priorities []string, trans ...transformer.Transformer) (interface{}, error)
	All(ctx context.Context, param helper.QueryParams, priorities []string, trans ...transformer.Transformer) (interface{}, error)
	Set(ctx context.Context, model interface{}) (int64, error)
	Query(ctx context.Context, params []helper.QueryParams, priorities []string, trans ...transformer.Transformer) (interface{}, error)
	Count(ctx context.Context, params []helper.QueryParams, priorities []string) (total int, err error)
	Insert(ctx context.Context, object interface{}) (int64, error)
	Update(ctx context.Context, model interface{}, params helper.QueryParams, priorities []string) error
	Upsert(ctx context.Context, model interface{}) (int64, error)
	Delete(ctx context.Context, params helper.QueryParams, priorities ...string) error
}
