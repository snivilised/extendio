package lo

// Entry defines a key/value pairs.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Tuple2 is a group of 2 elements (pair).
type Tuple2[A any, B any] struct {
	A A
	B B
}

// Unpack returns values contained in tuple.
func (t Tuple2[A, B]) Unpack() (A, B) { //nolint:gocritic // foo
	return t.A, t.B
}

// Tuple3 is a group of 3 elements.
type Tuple3[A any, B any, C any] struct {
	A A
	B B
	C C
}

// Unpack returns values contained in tuple.
func (t Tuple3[A, B, C]) Unpack() (A, B, C) { //nolint:gocritic // foo
	return t.A, t.B, t.C
}

// Tuple4 is a group of 4 elements.
type Tuple4[A any, B any, C any, D any] struct {
	A A
	B B
	C C
	D D
}

// Unpack returns values contained in tuple.
func (t Tuple4[A, B, C, D]) Unpack() (A, B, C, D) { //nolint:gocritic // foo
	return t.A, t.B, t.C, t.D
}

// Tuple5 is a group of 5 elements.
type Tuple5[A any, B any, C any, D any, E any] struct {
	A A
	B B
	C C
	D D
	E E
}

// Unpack returns values contained in tuple.
func (t Tuple5[A, B, C, D, E]) Unpack() (A, B, C, D, E) { //nolint:gocritic // foo
	return t.A, t.B, t.C, t.D, t.E
}

// Tuple6 is a group of 6 elements.
type Tuple6[A any, B any, C any, D any, E any, F any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
}

// Unpack returns values contained in tuple.
func (t Tuple6[A, B, C, D, E, F]) Unpack() (A, B, C, D, E, F) { //nolint:gocritic // foo
	return t.A, t.B, t.C, t.D, t.E, t.F
}

// Tuple7 is a group of 7 elements.
type Tuple7[A any, B any, C any, D any, E any, F any, G any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
	G G
}

// Unpack returns values contained in tuple.
func (t Tuple7[A, B, C, D, E, F, G]) Unpack() (A, B, C, D, E, F, G) { //nolint:gocritic // foo
	return t.A, t.B, t.C, t.D, t.E, t.F, t.G
}

// Tuple8 is a group of 8 elements.
type Tuple8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
	G G
	H H
}

// Unpack returns values contained in tuple.
func (t Tuple8[A, B, C, D, E, F, G, H]) Unpack() (A, B, C, D, E, F, G, H) { //nolint:gocritic // foo
	return t.A, t.B, t.C, t.D, t.E, t.F, t.G, t.H
}

// Tuple9 is a group of 9 elements.
type Tuple9[A any, B any, C any, D any, E any, F any, G any, H any, I any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
	G G
	H H
	I I
}

// Unpack returns values contained in tuple.
func (t Tuple9[A, B, C, D, E, F, G, H, I]) Unpack() (A, B, C, D, E, F, G, H, I) { //nolint:gocritic // foo
	return t.A, t.B, t.C, t.D, t.E, t.F, t.G, t.H, t.I
}
