package web

import (
	"bytes"
	"hash/adler32"
	"net"
	"net/http"
	"strings"
)

type iphasher struct{}

// ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

var additionalPrivateRanges = []ipRange{
	{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// inRange - check to see if a given ip address is within a range given
func (iph *iphasher) inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func (iph *iphasher) isPrivateSubnet(ipAddress net.IP) bool {
	// net.IP has some info about private IPs, but not all ranges included
	if ipAddress.IsPrivate() {
		return true
	}

	// check for additional private ranges
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		for _, r := range additionalPrivateRanges {
			if iph.inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

// match performs an actual IP matching
func (iph *iphasher) match(request *http.Request) string {
	for _, header := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(request.Header.Get(header), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || iph.isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	return ""
}

// GetHash returns hash of an IP address
func (iph *iphasher) GetHash(request *http.Request) uint32 {
	ip := iph.match(request)

	return adler32.Checksum([]byte(ip))
}
