package result

type Result[T any] struct {
	value T
	error error
}

func (r *Result[T]) Ok(value T) {
	r.value = value
}

func (r *Result[T]) Err(error error) {
	r.error = error
}
