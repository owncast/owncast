//go:build windows
// +build windows

package chat

func getMaximumConcurrentConnectionLimit() uint64 {
	// The maximum limit I can find for windows is 16,777,216
	// (essentially unlimited, but add the 0.7 multiplier as well to be
	// consistent with other systems)
	// https://docs.microsoft.com/en-gb/archive/blogs/markrussinovich/pushing-the-limits-of-windows-handles
	return (16777216 * 7) / 10
}
