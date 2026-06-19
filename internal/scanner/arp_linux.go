//go:build linux

package scanner

import (
	"bufio"
	"os"
	"strings"
)

func lookupMACFromARP(ip string) string {
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return ""
	}

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 4 && fields[0] == ip {
			return strings.ToUpper(fields[3])
		}
	}

	return ""
}
