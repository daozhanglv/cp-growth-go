package copypasta

/*
Berlekamp-Massey Algorithm (BM算法 / 线性递推)
https://en.wikipedia.org/wiki/Berlekamp%E2%80%93Massey_algorithm
https://oi-wiki.org/math/berlekamp-massey/

Given a sequence s[0..n-1] over Z/pZ (p prime), finds the shortest linear
recurrence: s[i] = rec[0]*s[i-1] + rec[1]*s[i-2] + ... + rec[L-1]*s[i-L].

See also: math_matrix.go (berlekampMassey) for the fixed-modulus int variant.

Typical use in contests:
  1. Compute enough terms with brute force or DP.
  2. Call berlekampMasseyMod to find the recurrence.
  3. Evaluate s[n] for large n using matrix exponentiation (O(L³ log n))
     or Cayley-Hamilton (O(L² log n + L log L log n)).

Time: O(n·L) where L = length of recurrence found.

Related problems:
https://codeforces.com/problemset/problem/1511/G  3100
https://www.luogu.com.cn/problem/P5487
https://loj.ac/p/6466
*/

// berlekampMasseyMod finds the shortest linear recurrence for s modulo mod (prime).
// Returns rec where s[i] = sum(rec[j]*s[i-1-j], j=0..L-1) (mod mod).
// Safe for negative or unreduced (> mod) int64 values in s.
func berlekampMasseyMod(s []int64, mod int64) []int64 {
	n := len(s)
	a := make([]int64, n)
	for i, v := range s {
		a[i] = (v%mod + mod) % mod
	}
	// C is the current LFSR polynomial: C[0]=1, recurrence encoded in C[1..L]
	// a[i] + C[1]*a[i-1] + ... + C[L]*a[i-L] = 0 (mod mod)
	C := []int64{1}
	B := []int64{1} // previous best C
	l := 0          // current LFSR length
	x := 1          // steps since last update
	b := int64(1)   // last non-zero discrepancy

	for i := 0; i < n; i++ {
		d := a[i]
		for j := 1; j <= l; j++ {
			d = (d + C[j]*a[i-j]) % mod
		}
		if d == 0 {
			x++
			continue
		}
		T := make([]int64, len(C))
		copy(T, C)
		// coef = d * b^{-1} mod p
		coef := d * bmPow(b, mod-2, mod) % mod
		// C = C - coef * x^shift * B
		shift := x
		for len(C) < len(B)+shift {
			C = append(C, 0)
		}
		for j, v := range B {
			C[j+shift] = (C[j+shift] - coef*v%mod + mod) % mod
		}
		if 2*l <= i {
			l = i + 1 - l
			B = T
			b = d
			x = 1
		} else {
			x++
		}
	}
	// Convert LFSR polynomial to recurrence coefficients:
	// s[i] = -C[1]*s[i-1] - ... - C[L]*s[i-L]   =>   rec[j] = (mod - C[j+1]) % mod
	rec := make([]int64, len(C)-1)
	for j, v := range C[1:] {
		rec[j] = (mod - v) % mod
	}
	return rec
}

// kitamasaMod evaluates the k-th term (0-indexed) of the linear recurrence
// s[i] = sum(rec[j]*s[i-1-j], j=0..L-1) (mod mod) in O(L^2 log k) time.
func kitamasaMod(rec, init []int64, k, mod int64) int64 {
	if k < int64(len(init)) {
		return (init[k]%mod + mod) % mod
	}
	L := len(rec)
	if L == 0 {
		return 0
	}
	normRec := make([]int64, L)
	for i, v := range rec {
		normRec[i] = (v%mod + mod) % mod
	}
	if L == 1 {
		a0 := (init[0]%mod + mod) % mod
		return a0 * bmPow(normRec[0], k, mod) % mod
	}

	// Multiply two polynomials of degree < L modulo x^L - sum_{j=0..L-1} rec[j]*x^{L-1-j}.
	mulPoly := func(p, q []int64) []int64 {
		prod := make([]int64, 2*L-1)
		for i, vp := range p {
			if vp == 0 {
				continue
			}
			for j, vq := range q {
				prod[i+j] = (prod[i+j] + vp*vq) % mod
			}
		}
		for i := 2*L - 2; i >= L; i-- {
			if prod[i] == 0 {
				continue
			}
			c := prod[i]
			for j := 0; j < L; j++ {
				prod[i-1-j] = (prod[i-1-j] + c*normRec[j]) % mod
			}
		}
		return prod[:L]
	}

	res := make([]int64, L)
	res[0] = 1
	base := make([]int64, L)
	base[1] = 1
	for exp := k; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = mulPoly(res, base)
		}
		base = mulPoly(base, base)
	}

	var ans int64
	for i := 0; i < L; i++ {
		ai := (init[i]%mod + mod) % mod
		ans = (ans + res[i]*ai) % mod
	}
	return ans
}

// guessNthMod finds the shortest recurrence for s modulo prime mod and evaluates s[k] (0-indexed).
func guessNthMod(s []int64, k, mod int64) int64 {
	if len(s) == 0 {
		return 0
	}
	if k < int64(len(s)) {
		return (s[k]%mod + mod) % mod
	}
	rec := berlekampMasseyMod(s, mod)
	return kitamasaMod(rec, s, k, mod)
}

func bmPow(base, exp, mod int64) int64 {
	result := int64(1)
	base = (base%mod + mod) % mod
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return result
}
