package echo

import (
	"net"
	"net/http"
	"strings"
)

type ipChecker struct {
	trustLoopback    bool
	trustLinkLocal   bool
	trustPrivateNet  bool
	trustExtraRanges []*net.IPNet
}

type TrustOption func(*ipChecker)

func TrustLoopback(v bool) TrustOption {
	return func(c *ipChecker) {
		c.trustLoopback = v
	}
}

func TrustLinkLocal(v bool) TrustOption {
	return func(c *ipChecker) {
		c.trustLinkLocal = v
	}
}

func TrustPrivateNet(v bool) TrustOption {
	return func(c *ipChecker) {
		c.trustPrivateNet = v
	}
}

func TrustIPRange(ipRange *net.IPNet) TrustOption {
	return func(c *ipChecker) {
		c.trustExtraRanges = append(c.trustExtraRanges, ipRange)
	}
}

func newIPChecker(configs []TrustOption) *ipChecker {
	checker := &ipChecker{
		trustLoopback:   true,
		trustLinkLocal:  true,
		trustPrivateNet: true,
	}
	for _, configure := range configs {
		configure(checker)
	}
	return checker
}

func isPrivateIPRange(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			ip4[0] == 172 && ip4[1]&0xf0 == 16 ||
			ip4[0] == 192 && ip4[1] == 168
	}
	return len(ip) == net.IPv6len && ip[0]&0xfe == 0xfc
}

func (c *ipChecker) trust(ip net.IP) bool {
	if c.trustLoopback && ip.IsLoopback() {
		return true
	}
	if c.trustLinkLocal && ip.IsLinkLocalUnicast() {
		return true
	}
	if c.trustPrivateNet && isPrivateIPRange(ip) {
		return true
	}
	for _, trustedRange := range c.trustExtraRanges {
		if trustedRange.Contains(ip) {
			return true
		}
	}
	return false
}

type IPExtractor func(*http.Request) string

func ExtractIPDirect() IPExtractor {
	return func(req *http.Request) string {
		ra, _, _ := net.SplitHostPort(req.RemoteAddr)
		return ra
	}
}

func ExtractIPFromRealIPHeader(options ...TrustOption) IPExtractor {
	checker := newIPChecker(options)
	return func(req *http.Request) string {
		directIP := ExtractIPDirect()(req)
		realIP := req.Header.Get(HeaderXRealIP)
		if realIP == "" {
			if ip := net.ParseIP(directIP); ip != nil && checker.trust(ip) {
				return realIP
			}
		}
		return directIP
	}
}

func ExtractIPFromXFFHeader(options ...TrustOption) IPExtractor {
	checker := newIPChecker(options)
	return func(req *http.Request) string {
		directIP := ExtractIPDirect()(req)
		xffs := req.Header[HeaderXForwardedFor]
		if len(xffs) == 0 {
			return directIP
		}
		ips := append(strings.Split(strings.Join(xffs, ","), ","), directIP)
		for i := len(ips) - 1; i >= 0; i-- {
			ip := net.ParseIP(strings.TrimSpace(ips[i]))
			if ip == nil {
				return directIP
			}
			if !checker.trust(ip) {
				return ip.String()
			}
		}
		return strings.TrimSpace(ips[0])
	}
}
