package mathx

// QMI 快速幂求 a ^ k % p
func QMI(a, k, p int) int {
	res := 1
	for k > 0 {
		if k&1 == 1 {
			res = res * a % p
		}
		a = a * a % p
		k >>= 1
	}
	return res
}
