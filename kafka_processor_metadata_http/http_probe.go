package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type HttpProbe struct {
	Client *http.Client
	Secure bool
}

type HttpProbeResult struct {
	Header http.Header
	StatusCode int
}

func (hp *HttpProbe) Get(ip string, port int) (HttpProbeResult, error) {
	var protocol string
	if hp.Secure {
		protocol = "https"
	} else {
		protocol = "http"
	}

	url := fmt.Sprintf("%s://%s:%d", protocol, ip, port)
	log.Println("running probe on ", url)

	resp, err := hp.Client.Get(url)
	if err != nil {
		return HttpProbeResult{}, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return HttpProbeResult{Header:resp.Header, StatusCode:resp.StatusCode}, err
}

type HttpProbeOptions struct {
	Timeout time.Duration
	Secure bool
}

func NewHttpProbe(options HttpProbeOptions) HttpProbe {
	var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   options.Timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}

	var re = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	client := &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       options.Timeout,
	}

	return HttpProbe{
		Client: client,
		Secure: options.Secure,
	}
}
