package storage

func MapSlice[S any, D any](src []S, f func(S) D) []D {
	dst := make([]D, len(src))
	for i, s := range src {
		dst[i] = f(s)
	}
	return dst
}
