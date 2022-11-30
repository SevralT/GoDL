package proxy

import (
	"net/http"
	"net/url"

	"github.com/SevralT/GoDL/vars"
)

func GetProxyClient() *http.Client {
	// Check if proxy server is setted
	if vars.ProxyServer != "" {
		// Check if proxy port is setted
		if vars.ProxyPort != "" {
			// Create proxy url
			proxyUrl, _ := url.Parse("http://" + vars.ProxyServer + ":" + vars.ProxyPort)
			// Create client
			client := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyUrl),
				},
			}
			return client
		} else {
			// Create proxy url
			proxyUrl, _ := url.Parse("http://" + vars.ProxyServer + ":8080")
			// Create client
			client := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyUrl),
				},
			}
			return client
		}
	} else {
		// Create client
		client := &http.Client{}
		return client
	}
}
