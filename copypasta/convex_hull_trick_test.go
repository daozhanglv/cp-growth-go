package copypasta

import (
	"math/rand"
	"testing"
)

func TestCHTDequeBasic(t *testing.T) {
	// lines with non-decreasing k: y=-x+20, y=x+10, y=3x+4, y=5x
	c := &chtDeque{}
	c.AddLine(-1, 20)
	c.AddLine(1, 10)
	c.AddLine(3, 4)
	c.AddLine(5, 0)

	// queries in non-decreasing x order
	cases := [][2]int{{0, 20}, {2, 18}, {5, 25}, {10, 50}, {15, 75}}
	for _, tc := range cases {
		if got := c.QueryMax(tc[0]); got != tc[1] {
			t.Errorf("QueryMax(%d): want %d, got %d", tc[0], tc[1], got)
		}
	}
}

func TestCHTDequeInterleavedAndEqualSlopes(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	type line struct{ k, b int }
	var lines []line
	cd := &chtDeque{}

	k := -500
	x := -500
	for step := 0; step < 500; step++ {
		// Advance k non-decreasingly (allow equal slopes).
		k += rng.Intn(3)
		b := rng.Intn(2000) - 1000
		lines = append(lines, line{k, b})
		cd.AddLine(k, b)

		// Interleaved non-decreasing x queries.
		for q := 0; q < 3; q++ {
			x += rng.Intn(3)
			brute := -1 << 60
			for _, l := range lines {
				if v := l.k*x + l.b; v > brute {
					brute = v
				}
			}
			if got := cd.QueryMax(x); got != brute {
				t.Fatalf("step %d x=%d: brute=%d chtDeque=%d", step, x, brute, got)
			}
		}
	}
}

func TestCHTStackArbitraryQuery(t *testing.T) {
	c := &chtStack{}
	c.AddLine(5, 0)
	c.AddLine(3, 4)
	c.AddLine(1, 10)
	c.AddLine(-1, 20)

	cases := map[int]int{
		0:  20,
		2:  18,
		5:  25,
		10: 50,
		15: 75,
		3:  17, // 5*3=15, 3*3+4=13, 1*3+10=13, -1*3+20=17
	}
	for x, want := range cases {
		if got := c.QueryMax(x); got != want {
			t.Errorf("QueryMax(%d): want %d, got %d", x, want, got)
		}
	}
}

func TestCHTConsistencyWithBrute(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	type line struct{ k, b int }
	var lines []line

	cs := &chtStack{}
	k := 100
	// add lines with non-increasing slope (including equal slopes)
	for i := 0; i < 80; i++ {
		k -= rng.Intn(3)
		b := rng.Intn(200) - 100
		lines = append(lines, line{k, b})
		cs.AddLine(k, b)
	}

	// verify against brute force for random queries
	for q := 0; q < 200; q++ {
		x := rng.Intn(200) - 100
		brute := -1 << 60
		for _, l := range lines {
			if v := l.k*x + l.b; v > brute {
				brute = v
			}
		}
		if got := cs.QueryMax(x); got != brute {
			t.Fatalf("x=%d: brute=%d cht=%d", x, brute, got)
		}
	}
}

func TestCHTLargeCoordinatesNoOverflow(t *testing.T) {
	// Slopes ~ 1e9 and intercepts ~ 1e12 would overflow int64 in naive (db * dk).
	const scaleK = 1_000_000_000
	const scaleB = 1_000_000_000_000
	lines := []chtLine{
		{1 * scaleK, 3 * scaleB},
		{2 * scaleK, 2 * scaleB},
		{3 * scaleK, 0},
	}

	cd := &chtDeque{}
	for _, l := range lines {
		cd.AddLine(l.k, l.b)
	}
	for _, x := range []int{0, 500, 1000, 1500, 3000} {
		want := lines[0].eval(x)
		for _, l := range lines[1:] {
			if v := l.eval(x); v > want {
				want = v
			}
		}
		if got := cd.QueryMax(x); got != want {
			t.Fatalf("chtDeque large x=%d: want %d got %d", x, want, got)
		}
	}

	cs := &chtStack{}
	for i := len(lines) - 1; i >= 0; i-- {
		cs.AddLine(lines[i].k, lines[i].b)
	}
	for _, x := range []int{3000, 0, 1500, 500, 1000} {
		want := lines[0].eval(x)
		for _, l := range lines[1:] {
			if v := l.eval(x); v > want {
				want = v
			}
		}
		if got := cs.QueryMax(x); got != want {
			t.Fatalf("chtStack large x=%d: want %d got %d", x, want, got)
		}
	}
}
