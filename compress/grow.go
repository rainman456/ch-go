package compress

// growBuf returns b with length n, reusing the backing array when the capacity
// is sufficient. The retained bytes are not reset: callers must overwrite or
// explicitly zero every byte they read from the result afterwards. Otherwise it
// grows the buffer the way this package always has, with a fresh zeroed region.
//
// Resetting the buffer on every block would zero bytes that are rewritten or
// discarded immediately afterwards, and the zero fill was measurably hot in the
// package benchmarks. See https://github.com/ClickHouse/ch-go/issues/1218.
func growBuf(b []byte, n int) []byte {
	if cap(b) >= n {
		return b[:n]
	}
	return append(b[:0], make([]byte, n)...)
}
