// Package opt provides a simple, generic option type.
package opt

// Opt is an option type.
//
// It contains either contains a value or does not.
type Opt[T any] struct {
	// INVARIANT: !isSome => value == zeroValue[T]()
	// ...so that: (!o1.IsSome() && !o2.IsSome()) => (o1 == o2)
	//
	// In English, this requires that if the isSome field is false, the value
	// field should be the zero value for the type T, so that the "==" operator
	// returns an accurate result for any two option values (as long as T is
	// comparable).

	// value is the wrapped value, if there is one.
	value T
	// isSome indicates whether the option contains a value.
	isSome bool
}

// None creates a new option that does not contain a value.
func None[T any]() Opt[T] { return Opt[T]{isSome: false} }

// Some creates a new option that contains the given value.
func Some[T any](v T) Opt[T] { return Opt[T]{value: v, isSome: true} }

// FromPtr converts a pointer to an option. A nil pointer becomes an option
// containing no value, and a non-nil pointer becomes an option containing the
// value pointed to by the pointer.
func FromPtr[T any](p *T) Opt[T] {
	if p == nil {
		return None[T]()
	}

	return Some(*p)
}

// FromValOk converts a value and a bool to an option. When ok is false the
// result is an option containing no value, and when ok is true the result
// contains v.
func FromValOk[T any](v T, ok bool) Opt[T] {
	if !ok {
		return None[T]()
	}

	return Some(v)
}

// FromValErr converts a value and an error to an option. When err is non-nil
// the result is an option containing no value, and when err is nil the result
// contains v.
func FromValErr[T any](v T, err error) Opt[T] {
	if err != nil {
		return None[T]()
	}

	return Some(v)
}

// Get returns the wrapped value and true if the option contains a value.
// Otherwise, it returns the zero value for T and false.
func (o Opt[T]) Get() (v T, ok bool) {
	if !o.isSome {
		var z T
		return z, false
	}

	return o.value, true
}

// IsSome returns whether the option contains a value.
func (o Opt[T]) IsSome() bool { return o.isSome }

// Unwrap returns the wrapped value if the option contains a value. Otherwise,
// it panics.
func (o Opt[T]) Unwrap() T {
	if !o.isSome {
		panic("None unwrapped")
	}

	return o.value
}

// UnwrapOr returns the wrapped value if the option contains a value.
// Otherwise, it returns v.
func (o Opt[T]) UnwrapOr(v T) T {
	if !o.isSome {
		return v
	}

	return o.value
}

// UnwrapOrElse returns the wrapped value if the option contains a value,
// otherwise, it returns the result of f.
func (o Opt[T]) UnwrapOrElse(f func() T) T {
	if !o.isSome {
		return f()
	}

	return o.value
}

// Take returns the value of o, replacing it with an option that contains no
// value.
func (o *Opt[T]) Take() Opt[T] {
	if !o.isSome {
		return None[T]()
	}

	res := *o
	*o = None[T]()
	return res
}

// ToPtr converts the option to a pointer. An option that contains no value
// becomes nil, and an option that does contain a value becomes a pointer to
// that value.
func (o Opt[T]) ToPtr() *T {
	if !o.isSome {
		return nil
	}

	return &o.value
}

// Map returns a new option, which is the result of applying f to the value
// wrapped by o, if it contains a value. If o does not contain a value, f is
// not evaluated, and the returned option does not contain a value either.
func Map[T, U any](o Opt[T], f func(T) U) Opt[U] {
	if !o.isSome {
		return None[U]()
	}

	return Some(f(o.value))
}

// Match returns the result of evaluating some with the value wrapped by o if
// o contains a value, or the result of evaluating none if it does not.
func Match[T, U any](o Opt[T], some func(T) U, none func() U) U {
	if o.isSome {
		return some(o.value)
	}

	return none()
}
