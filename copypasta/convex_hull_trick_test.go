package copypasta

import (
	"math/rand"
	"testing"
)

func TestCHTDequeBasic(t *testing.T) {
	// lines with decreasing k: y=5x, y=3x+4, y=x+10, y=-x+20
	c := &chtDeque{}
	c.AddLine(5, 0)
	c.AddLine(3, 4)
	c.AddLine(1, 10)
	c.AddLine(-1, 20)

	// queries in increasing x order
	cases := [][2]int{{0, 20}, {2, 18}, {5, 25}, {10, 50}, {15, 75}}
	// manually computed maximums:
	// x=0:  max(0, 4, 10, 20) = 20  (line -x+20)
	// x=2:  max(10, 10, 12, 18) = 18  (line -x+20... wait let me recalculate)
	// Actually: at x=2: 5*2=10, 3*2+4=10, 1*2+10=12, -1*2+20=18 → max=18
	// at x=5: 25, 19, 15, 15 → max=25
	// at x=10: 50, 34, 20, 10 → max=50
	// at x=15: 75, 49, 25, 5 → max=75
	expected := []int{20, 18, 25, 50, 75}
	for i, tc := range cases {
		if got := c.QueryMax(tc[0]); got != expected[i] {
			t.Errorf("QueryMax(%d): want %d, got %d", tc[0], expected[i], got)
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
	n := 50
	// add lines with decreasing slope
	for i := 0; i < n; i++ {
		k := n - i
		b := rng.Intn(200) - 100
		lines = append(lines, line{k, b})
		cs.AddLine(k, b)
	}

	// verify against brute force for random queries
	for q := 0; q < 100; q++ {
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
