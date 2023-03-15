package opt

import (
	"reflect"
	"testing"
)

func TestNone(t *testing.T) {
	assertEqual(t, Opt[string]{}, None[string]())
}

func TestSome(t *testing.T) {
	assertEqual(t, Opt[int]{value: 7, isSome: true}, Some(7))
}

func TestFromPtr(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if FromPtr[complex64](nil).IsSome() {
			t.Fatal("shouldn't be some")
		}
	})

	t.Run("non-nil", func(t *testing.T) {
		type foo struct {
			a, b uint
		}

		expectedValue := foo{a: 7, b: 9}

		expected := Opt[foo]{value: expectedValue, isSome: true}

		assertEqual(t, expected, FromPtr(&expectedValue))
	})
}

func TestOpt_Get(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		actual, actualOk := None[int]().Get()
		if actualOk {
			t.Fatal("should've been false")
		}

		var z int
		if actual != z {
			t.Fatal("should've been zero value")
		}
	})

	t.Run("some", func(t *testing.T) {
		const expected = 367

		actual, actualOk := Some(expected).Get()
		if !actualOk {
			t.Fatal("should've been true")
		}

		assertEqual(t, expected, actual)
	})
}

func TestOpt_IsSome(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		if None[int16]().IsSome() {
			t.Fatal("shouldn't be some")
		}
	})

	t.Run("some", func(t *testing.T) {
		if !Some(int16(9)).IsSome() {
			t.Fatal("shouldn't be none")
		}
	})
}

func TestOpt_Unwrap(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		defer func() {
			assertEqual(t, "None unwrapped", recover())
		}()

		None[complex128]().Unwrap()

		t.Fatal("should've panicked")
	})

	t.Run("some", func(t *testing.T) {
		assertEqual(t, 9.79, Some(9.79).Unwrap())
	})
}

func TestOpt_UnwrapOr(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		assertEqual(t, 7, None[int]().UnwrapOr(7))
	})

	t.Run("some", func(t *testing.T) {
		assertEqual(t, -976.165, Some(-976.165).UnwrapOr(397.319))
	})
}

func TestOpt_UnwrapOrElse(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		const expected = 7

		elseRun := false
		actual := None[int]().UnwrapOrElse(func() int {
			elseRun = true
			return expected
		})

		if !elseRun {
			t.Fatal("else not run")
		}

		assertEqual(t, expected, actual)
	})

	t.Run("some", func(t *testing.T) {
		const expected = -976.165

		actual := Some(expected).UnwrapOrElse(func() float64 {
			t.Fatal("else should not have been run")
			return 0
		})

		assertEqual(t, expected, actual)
	})
}

func TestOpt_Take(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		o := None[bool]()
		assertEqual(t, None[bool](), o.Take())
		assertEqual(t, None[bool](), o)
	})

	t.Run("some", func(t *testing.T) {
		o := Some(true)
		assertEqual(t, Some(true), o.Take())
		assertEqual(t, None[bool](), o)
	})
}

func TestOpt_ToPtr(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		assertEqual(t, (*[]string)(nil), None[[]string]().ToPtr())
	})

	t.Run("some", func(t *testing.T) {
		expected := [3]int{8, 3, 9}

		assertEqual(t, &expected, Some(expected).ToPtr())
	})
}

func TestMap(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		assertEqual(t, None[uint](), Map(None[int](), func(i int) uint {
			t.Fatal("map func should not have been run")
			return 0
		}))
	})

	t.Run("some", func(t *testing.T) {
		funcRun := false
		assertEqual(t, Some(11), Map(Some("hello world"), func(s string) int {
			funcRun = true
			assertEqual(t, "hello world", s)
			return len(s)
		}))

		if !funcRun {
			t.Fatal("map func should've been run")
		}
	})
}

func TestMatch(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		const expected = 79

		noneRun := false
		actual := Match(None[bool](),
			func(b bool) int {
				t.Fatal("some branch should not have been run")
				return -7
			},
			func() int {
				noneRun = true
				return expected
			},
		)

		if !noneRun {
			t.Fatal("none branch not run")
		}
		assertEqual(t, expected, actual)
	})

	t.Run("some", func(t *testing.T) {
		expected := [7]byte{1, 2, 3, 4, 5, 6, 7}

		someRun := false
		actual := Match(Some([]string{"foo", "bar"}),
			func(v []string) [7]byte {
				someRun = true
				assertEqual(t, []string{"foo", "bar"}, v)
				return expected
			},
			func() [7]byte {
				t.Fatal("none branch should not have been run")
				return [7]byte{}
			},
		)

		if !someRun {
			t.Fatal("some branch not run")
		}
		assertEqual(t, expected, actual)
	})
}

func assertEqual(t *testing.T, expected, actual any) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %#v, got: %#v", expected, actual)
	}
}
