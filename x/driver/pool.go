package driver

import "fmt"

type Pool[K comparable, D any] struct {
	drivers     map[K]*Factory[K, D]
	driverCache map[K]D
}

func NewPool[K comparable, D any](drivers Factories[K, D]) *Pool[K, D] {
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

func (p *Pool[K, D]) All() ([]D, error) {
	var r []D
	for k := range p.drivers {
		if d, err := p.Resolve(k); err != nil {
			if f, ok := p.drivers[k]; ok && f.Optional {
				// if the driver is optional, we skip it
				continue
			}

			return nil, fmt.Errorf("failed to resolve driver %v: %w", k, err)
		} else {
			r = append(r, d)
		}
	}

	return r, nil
}
