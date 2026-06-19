//go:build !linux

package scanner

func lookupMACFromARP(ip string) string {
	return ""
}
