package driver

import (
	"go.uber.org/fx"
)

type Factory[K comparable, D any] struct {
	Provides K
	Create   func() (D, error)
	Optional bool
}

type Factories[K comparable, D any] []*Factory[K, D]

func NewFactory[K comparable, D any](name K, create func() (D, error)) FactoryResult[K, D] {
	return FactoryResult[K, D]{
		Factory: &Factory[K, D]{
			Provides: name,
			Create:   create,
		},
	}
}

func NewOptionalFactory[K comparable, D any](name K, create func() (D, error)) FactoryResult[K, D] {
	return FactoryResult[K, D]{
		Factory: &Factory[K, D]{
			Provides: name,
			Create:   create,
			Optional: true,
		},
	}
}

type FactoryResult[K comparable, D any] struct {
	fx.Out

	Factory *Factory[K, D] `group:"drivers"`
}
