package main

import (
	"net/http"
)

type JsonOutput struct {
	Ip string
	Port int
	Status int
	Protocol string
	Country string
	Headers http.Header
}
