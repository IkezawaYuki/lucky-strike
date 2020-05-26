package middleware

import (
	echo "github.com/IkezawaYuki/lucky-strike"
	"math/rand"
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

var ()
