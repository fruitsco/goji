package x

func MergeMap[M ~map[K]V, K comparable, V any](maps ...M) M {
	fullCap := 0
	for _, m := range maps {
		fullCap += len(m)
	}

	merged := make(M, fullCap)
	for _, m := range maps {
		for key, val := range m {
			merged[key] = val
		}
	}

	return merged
}
