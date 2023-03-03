package optional

// アロケーションやGC管理、nilチェックのオーバーヘッドを防ぐためにポインタではなくbool値で管理する

type Option[T any] struct {
	v      T
	isSome bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{
		v:      v,
		isSome: true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func FromNillable[T any](v *T) Option[T] {
	if v == nil {
		return None[T]()
	}
	return Some(*v)
}

func PtrFromNillable[T any](ptr *T) Option[*T] {
	if ptr == nil {
		return None[*T]()
	}
	return Some(ptr)
}

func (o Option[T]) IsNone() bool {
	return !o.isSome
}

func (o Option[T]) IsSome() bool {
	return o.isSome
}

// Unwrap returns the value regardless of Some/None status.
// If the Option value is Some, this method returns the actual value.
// On the other hand, if the Option value is None, this method returns the *default* value according to the type.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		var defaultValue T
		return defaultValue
	}
	return o.v
}

// Take takes the contained value in Option.
// If Option value is Some, this returns the value that is contained in Option.
// On the other hand, this returns an ErrNoneValueTaken as the second return value.
func (o Option[T]) Take() (T, bool) {
	if o.IsNone() {
		var defaultValue T
		return defaultValue, false
	}
	return o.v, true
}

// TakeOr returns the actual value if the Option has a value.
// On the other hand, this returns fallbackValue.
func (o Option[T]) TakeOr(fallbackValue T) T {
	if o.IsNone() {
		return fallbackValue
	}
	return o.v
}

// TakeOrElse returns the actual value if the Option has a value.
// On the other hand, this executes fallbackFunc and returns the result value of that function.
func (o Option[T]) TakeOrElse(fallbackFunc func() T) T {
	if o.IsNone() {
		return fallbackFunc()
	}
	return o.v
}
