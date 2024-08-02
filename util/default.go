package util

func WithDefault[T any](in *T, def T) T {
	if in == nil {
		return def
	}

	return *in
}

func WithDefaultFn[T any](in *T, def func() T) T {
	if in == nil {
		return def()
	}

	return *in
}
