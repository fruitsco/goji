package conf

func MergeDefaults[M ~map[string]V, V any](maps ...M) M {
	return MergeDefaultsNs("", maps...)
}

func MergeDefaultsNs[M ~map[string]V, V any](ns string, maps ...M) M {
	fullCap := 0
	for _, m := range maps {
		fullCap += len(m)
	}

	merged := make(M, fullCap)
	for _, m := range maps {
		for key, val := range m {
			newKey := key
			if ns != "" {
				newKey = ns + "." + key
			}
			merged[newKey] = val
		}
	}

	return merged
}
