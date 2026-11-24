package indieauth

import "crypto/rand"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randString returns n letters drawn from crypto/rand. Used for OAuth state
// and PKCE verifiers, so it must be unpredictable.
func randString(n int) string {
	b := make([]byte, n)
	buf := make([]byte, n)
	for i := 0; i < n; {
		_, _ = rand.Read(buf) // never errors; crashes the program on failure
		for _, c := range buf {
			// Rejection sampling keeps the 52 letters uniform.
			if c &= 63; c < byte(len(letterBytes)) && i < n {
				b[i] = letterBytes[c]
				i++
			}
		}
	}
	return string(b)
}
