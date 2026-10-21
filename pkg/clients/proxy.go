package clients

import (
	"net/http"
	"net/url"
)

// proxyForClientGo returns a proxy function that is compatible with all of
// client-go's transports, including the SPDY transport used by pod exec.
//
// client-go's SPDY transport supports the SOCKS5 protocol, but only accepts
// the "socks5" URL scheme. Tools such as sshuttle commonly advertise their
// local listener as "socks5h" to request proxy-side DNS resolution. The
// SOCKS5 dialer used by client-go already sends the original host name to the
// proxy, so these schemes are equivalent for this client.
func proxyForClientGo(proxy func(*http.Request) (*url.URL, error)) func(*http.Request) (*url.URL, error) {
	if proxy == nil {
		proxy = http.ProxyFromEnvironment
	}

	return func(request *http.Request) (*url.URL, error) {
		proxyURL, err := proxy(request)
		if err != nil || proxyURL == nil || proxyURL.Scheme != "socks5h" {
			return proxyURL, err
		}

		compatibleURL := *proxyURL
		compatibleURL.Scheme = "socks5"

		return &compatibleURL, nil
	}
}
