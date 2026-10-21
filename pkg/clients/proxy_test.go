package clients

import (
	"net/http"
	"net/url"
	"testing"
)

func TestProxyForClientGoConvertsSocks5h(t *testing.T) {
	proxyURL, err := url.Parse("socks5h://localhost:1080")
	if err != nil {
		t.Fatal(err)
	}

	proxy := proxyForClientGo(func(*http.Request) (*url.URL, error) {
		return proxyURL, nil
	})

	got, err := proxy(&http.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Scheme != "socks5" {
		t.Fatalf("proxy scheme = %q, want socks5", got.Scheme)
	}
	if proxyURL.Scheme != "socks5h" {
		t.Fatalf("input proxy scheme was changed to %q", proxyURL.Scheme)
	}
}

func TestProxyForClientGoUsesEnvironmentWhenProxyIsNil(t *testing.T) {
	t.Setenv("HTTPS_PROXY", "socks5h://localhost:1080")
	t.Setenv("NO_PROXY", "")

	proxy := proxyForClientGo(nil)
	got, err := proxy(&http.Request{URL: mustParseURL(t, "https://api.example.test")})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Scheme != "socks5" {
		t.Fatalf("proxy = %v, want socks5 proxy normalized to socks5", got)
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
