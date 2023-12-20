package driver

import (
	"go.uber.org/fx"
)

type Factory[K comparable, D any] struct {
	Provides K
	Create   func() (D, error)
}

func NewFactory[K comparable, D any](name K, create func() (D, error)) FactoryResult[K, D] {
	return FactoryResult[K, D]{
		Factory: &Factory[K, D]{
			Provides: name,
			Create:   create,
		},
	}
}

type FactoryResult[K comparable, D any] struct {
	fx.Out

	Factory *Factory[K, D] `group:"drivers"`
}
