package util

func Dereference[T any](p *T) T {
	var value T
	if p != nil {
		value = *p
	}
	return value
}
