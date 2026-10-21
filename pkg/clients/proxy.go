package clients

import (
	"net/http"
	"net/url"
)

const (
	socks5hScheme = "socks5h"
	socks5Scheme  = "socks5"
)

// proxyForClientGo returns a proxy function that is compatible with all of
// client-go's transports, including the SPDY transport used by pod exec.
//
// client-go's SPDY transport supports the SOCKS5 protocol, but only accepts
// the "socks5" URL scheme. Tools such as sshuttle commonly advertise their
// local listener as "socks5h" to request proxy-side DNS resolution. The
// SOCKS5 dialer used by client-go already sends the original host name to the
// proxy, so these schemes are equivalent for this client.
func proxyForClientGo(proxy func(*http.Request) (*url.URL, error), host string) func(*http.Request) (*url.URL, error) {
	if proxy == nil {
		proxy = http.ProxyFromEnvironment
	}

	request, err := http.NewRequest(http.MethodGet, host, nil)
	if err != nil {
		return nil
	}

	proxyURL, err := proxy(request)
	if err != nil || proxyURL == nil || proxyURL.Scheme != socks5hScheme {
		return nil
	}

	return func(request *http.Request) (*url.URL, error) {
		proxyURL, err := proxy(request)
		if err != nil || proxyURL == nil || proxyURL.Scheme != socks5hScheme {
			return proxyURL, err
		}

		compatibleURL := *proxyURL
		compatibleURL.Scheme = socks5Scheme

		return &compatibleURL, nil
	}
}
