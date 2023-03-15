# opt

Package opt provides a simple, generic option type.

## Why Not `*T`?

1. The semantics of `*T` can be unclear. Maybe this is a pointer type to allow
   shared access. Maybe it is a pointer to allow mutation. Maybe it is a
   pointer to deal with the idiosyncrasies of Go's encoding packages. If it's
   only a pointer because it may or may not be nil, the `Opt[T]` type
   communicates the value's significance more clearly.
2. `Opt[T]` encourages the consumer to handle the case where the option does
   not contain a value more explicitly. `*v` and especially `v.foo` are both
   much easier do without thinking than `o.Unwrap()`.
3. Because it may allow the wrapped value to be mutated unintentionally. (A
   similar problem can still arise with `Opt[T]` if the option contains a
   reference type, but it at least prevents direct reassignment of the value.)
4. Because it may force the value to be unnecessarily heap-allocated.

## Style Suggestions

These are just suggestions, exceptions may apply.

### `(v T, ok bool)` instead of `Opt[T]`

```go
// prefer:
func() (v T, ok bool) { ... }

// ...instead of:
func() Opt[T] { ... }
```

When returning a single value, the `(v T, ok bool)` pattern is generally more
idiomatic than returning a single option, unless the use-case for the function
depends upon the fact that it returns a value of the `Opt` type (such as the
`Map` function in this package).

### `if ... Get` instead of `if ... IsSome ... Unwrap`

```go
// prefer:
if v, ok := o.Get(); ok {
  _ = v
}

// ...or:
v, ok := o.Get()
if !ok {
  return
}

// ...instead of:
if o.IsSome() {
  _ = o.Unwrap()
}

// ...or:
if !o.IsSome() {
  return
}

_ = o.Unwrap()
```

This has no direct safety impact, but the first is more resilient to future
changes: if the `Unwrap` ends up getting surrounded by other stuff, the
relationship between the `IsSome` and the `Unwrap` may eventually become
unclear, leading to the introduction of a bug if the check is removed or the
use is relocated. The `Get` patterns are less risky because they follow the
common `(v T, ok bool)` idiom, meaning readers are less likely to separate the
check from the use.
