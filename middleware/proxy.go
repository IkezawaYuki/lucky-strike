package middleware

import (
	"fmt"
	echo "github.com/IkezawaYuki/lucky-strike"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sync"
)

type (
	ProxyConfig struct {
		Skipper      Skipper
		Balancer     ProxyBalancer
		Rewrite      map[string]string
		ContextKey   string
		Transport    http.RoundTripper
		rewriteRegex map[*regexp.Regexp]string
	}
	ProxyTarget struct {
		Name string
		URL  *url.URL
		Meta echo.Map
	}
	ProxyBalancer interface {
		AddTarget(*ProxyTarget) bool
		RemoveTarget(string) bool
		Next(echo.Context) *ProxyTarget
	}
	commonBalancer struct {
		targets []*ProxyTarget
		mutex   sync.RWMutex
	}
	randomBalancer struct {
		*commonBalancer
		random *rand.Rand
	}
	roundRobinBalancer struct {
		*commonBalancer
		i uint32
	}
)

var (
	DefaultProxyConfig = ProxyConfig{
		Skipper:    DefaultSkipper,
		ContextKey: "target",
	}
)

func proxyRaw(t *ProxyTarget, c echo.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, _, err := c.Response().Hijack()
		if err != nil {
			c.Set("_error", fmt.Sprintf("proxy raw, hijack error=%v, url=%s", t.URL, err))
			return
		}
		defer in.Close()

		out, err := net.Dial("tcp", t.URL.Host)
	})
}
