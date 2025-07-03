package nstore

type ConditionalFunc[T any] func(t T) bool

func (s *NStorage[T]) Query(cfunc ConditionalFunc[T], limit int) ([]T, int) {
	results := make([]T, limit)

	count := 0
	for _, item := range s.ListOfCache() {
		if cfunc(item) {
			results[count] = item
			count++
		}

		if count >= limit {
			break
		}
	}

	return results, count
}
