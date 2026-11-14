package copypasta

/*
Berlekamp-Massey Algorithm (BM算法 / 线性递推)
https://en.wikipedia.org/wiki/Berlekamp%E2%80%93Massey_algorithm
https://oi-wiki.org/math/berlekamp-massey/

Given a sequence s[0..n-1] over Z/pZ (p prime), finds the shortest linear
recurrence: s[i] = rec[0]*s[i-1] + rec[1]*s[i-2] + ... + rec[L-1]*s[i-L].

Typical use in contests:
  1. Compute enough terms with brute force or DP.
  2. Call berlekampMassey to find the recurrence.
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
func berlekampMasseyMod(s []int64, mod int64) []int64 {
	n := len(s)
	// C is the current LFSR polynomial: C[0]=1, recurrence encoded in C[1..L]
	// s[i] + C[1]*s[i-1] + ... + C[L]*s[i-L] = 0
	C := []int64{1}
	B := []int64{1} // previous best C
	l := 0          // current LFSR length
	x := 1          // steps since last update
	b := int64(1)   // last non-zero discrepancy

	for i := 0; i < n; i++ {
		// compute discrepancy d = s[i] + C[1]*s[i-1] + ...
		d := s[i] % mod
		for j := 1; j <= l; j++ {
			d = (d + C[j]*s[i-j]) % mod
		}
		d = (d%mod + mod) % mod
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

func bmPow(base, exp, mod int64) int64 {
	result := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return result
}
