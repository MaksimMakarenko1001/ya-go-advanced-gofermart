package pkg

func ValuesToList[K comparable, V any](m map[K]V) (l []V) {
	l = make([]V, 0, len(m))
	for _, v := range m {
		l = append(l, v)
	}
	return l
}

func KeysToList[K comparable, V any](m map[K]V) (l []K) {
	l = make([]K, 0, len(m))
	for k := range m {
		l = append(l, k)
	}
	return l
}

func Select[T, S any](l []T, f func(x T) S) (s []S) {
	s = make([]S, 0, len(l))
	for _, v := range l {
		s = append(s, f(v))
	}
	return s
}

// Creates a new pointer to the value.
func ToPtr[T any](value T) *T {
	return &value
}
