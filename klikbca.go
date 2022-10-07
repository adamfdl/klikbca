package klikbca

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gocolly/colly"
)

type klikBca struct {
	username      string
	password      string
	ipAddress     string
	colly         *colly.Collector
	proxyUrl      string
	delayDuration time.Duration
}

type option = func(k *klikBca)

func WithProxy(proxyUrl string) option {
	return func(klikBca *klikBca) {
		klikBca.proxyUrl = proxyUrl
	}
}

func WithDelay(dur time.Duration) option {
	return func(klikBca *klikBca) {
		klikBca.delayDuration = dur
	}
}

func NewKlikBca(userName, password string, opts ...option) *klikBca {

	klikBca := &klikBca{
		username: userName,
		password: password,
	}

	// Apply options
	for _, opt := range opts {
		opt(klikBca)
	}

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"

	// Set delay
	c.Limit(&colly.LimitRule{
		Delay:       klikBca.delayDuration,
		RandomDelay: klikBca.delayDuration,
	})

	// Skip TLS
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	// Set proxy
	if klikBca.proxyUrl != "" {
		c.SetProxy(klikBca.proxyUrl)
	}

	klikBca.colly = c

	return klikBca
}
