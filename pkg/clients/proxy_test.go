package clients

import (
	"net/http"
	"net/url"
	"testing"
)

const testAPIURL = "https://api.example.test"

func TestProxyForClientGoConvertsSocks5h(t *testing.T) {
	proxyURL, err := url.Parse(socks5hScheme + "://localhost:1080")
	if err != nil {
		t.Fatal(err)
	}

	proxy := proxyForClientGo(func(*http.Request) (*url.URL, error) {
		return proxyURL, nil
	}, testAPIURL)

	got, err := proxy(&http.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Scheme != socks5Scheme {
		t.Fatalf("proxy scheme = %q, want %s", got.Scheme, socks5Scheme)
	}
	if proxyURL.Scheme != socks5hScheme {
		t.Fatalf("input proxy scheme was changed to %q", proxyURL.Scheme)
	}
}

func TestProxyForClientGoUsesEnvironmentWhenProxyIsNil(t *testing.T) {
	t.Setenv("HTTPS_PROXY", socks5hScheme+"://localhost:1080")
	t.Setenv("NO_PROXY", "")

	proxy := proxyForClientGo(nil, testAPIURL)
	got, err := proxy(&http.Request{URL: mustParseURL(t, "https://api.example.test")})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Scheme != socks5Scheme {
		t.Fatalf("proxy = %v, want socks5 proxy normalized to socks5", got)
	}
}

func TestProxyForClientGoReturnsNilWhenConversionIsNotNeeded(t *testing.T) {
	proxyURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	proxy := proxyForClientGo(func(*http.Request) (*url.URL, error) {
		return proxyURL, nil
	}, testAPIURL)

	if proxy != nil {
		t.Fatal("proxy wrapper should be nil when SOCKS5 scheme conversion is not needed")
	}
}

func mustParseURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}
