package optionutil

func Configure[T any, O ~func(*T)](t *T, opts []O) *T {
	for _, opt := range opts {
		opt(t)
	}
	return t
}
