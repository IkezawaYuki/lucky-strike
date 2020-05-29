package middleware

import (
	"fmt"
	echo "github.com/IkezawaYuki/lucky-strike"
	"io"
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
		if err != nil {
			c.Set("_error", fmt.Sprintf("proxy raw, hijack error=%v, url=%s", t.URL, err))
			return
		}
		defer out.Close()

		err = r.Write(out)
		if err != nil {
			c.Set("_error", echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("proxy raw, dial error=%v, url=%s", t.URL, err)))
			return
		}
		defer out.Close()

		err = r.Write(out)
		if err != nil {
			c.Set("_error", echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("proxy raw, request header copy error=%v, url=%s", t.URL, err)))
			return
		}

		errCh := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err = io.Copy(dst, src)
			errCh <- err
		}
		go cp(out, in)
		go cp(in, out)
		err = <-errCh
		if err != nil && err != io.EOF {
			c.Set("_error", fmt.Errorf("proxy raw, copy body error=%v, url=%s", t.URL, err))
		}
	})
}

func NewRandomBalancer(targets []*ProxyTarget) ProxyBalancer {
	b := &randomBalancer{
		commonBalancer: new(commonBalancer),
	}
	b.targets = targets
	return b
}

func (b *commonBalancer) AddTarget(target *ProxyTarget) bool {
	for _, t := range b.targets {
		if t.Name == target.Name {
			return false
		}
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.targets = append(b.targets, target)
	return true
}