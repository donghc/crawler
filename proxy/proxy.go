package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type ProxyFunc func(r *http.Request) (*url.URL, error)

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

// GetProxy 取余算法实现轮询调度
func (r *roundRobinSwitcher) GetProxy(request *http.Request) (*url.URL, error) {
	index := atomic.AddUint32(&r.index, 1) - 1
	remainder := index % uint32(len(r.proxyURLs))
	u := r.proxyURLs[remainder]
	return u, nil
}

func RoundRobinProxySwitcher(proxyURLs ...string) (ProxyFunc, error) {
	if len(proxyURLs) < 1 {
		return nil, errors.New("proxy url list is empty")
	}
	urls := make([]*url.URL, len(proxyURLs))
	for i, u := range proxyURLs {
		parse, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parse
	}
	round := &roundRobinSwitcher{urls, 0}
	return round.GetProxy, nil
}
