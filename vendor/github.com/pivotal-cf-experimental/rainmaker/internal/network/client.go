package network

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var _client *http.Client

func GetClient(config Config) *http.Client {
	if _client == nil {
		_client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.SkipVerifySSL,
				},
				Proxy: http.ProxyFromEnvironment,
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		}
	}

	return _client
}
