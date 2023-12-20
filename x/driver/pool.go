package driver

import "fmt"

type Pool[K comparable, D any] struct {
	drivers     map[K]*Factory[K, D]
	driverCache map[K]D
}

func NewPool[K comparable, D any](drivers []*Factory[K, D]) *Pool[K, D] {
	driversPool := make(map[K]*Factory[K, D])
	driverCache := make(map[K]D)

	for _, factory := range drivers {
		driversPool[factory.Provides] = factory
	}

	return &Pool[K, D]{
		drivers:     driversPool,
		driverCache: driverCache,
	}
}

func (p *Pool[K, D]) Resolve(driverKey K) (D, error) {
	if driver, ok := p.driverCache[driverKey]; ok {
		return driver, nil
	}

	var d D

	if factory, ok := p.drivers[driverKey]; ok {
		driver, err := factory.Create()
		if err != nil {
			return d, err
		}

		p.driverCache[driverKey] = driver

		return driver, nil
	}

	return d, fmt.Errorf("driver %v not found", driverKey)
}
