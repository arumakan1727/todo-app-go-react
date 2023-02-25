package handler

type Handler[U any] struct {
	usecase U
}

func NewHandler[U any](usecase U) Handler[U] {
	return Handler[U]{
		usecase: usecase,
	}
}
