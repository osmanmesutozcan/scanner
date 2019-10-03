package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"io/ioutil"
	"log"
	"net"
	"scanner_workerpool"
	"time"
)

func main() {
	var port int
	var verbose bool
	var secure bool
	var outputraw bool
	var concurrency int
	timeout := 5000

	flag.BoolVar(&outputraw, "r", false, "output raw test instead of JSON")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&secure, "s", false, "secure")
	flag.IntVar(&port, "p", 80, "port")
	flag.IntVar(&concurrency, "c", 10, "concurrency")
	flag.Parse()

	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

	topic := flag.Arg(0)
	if len(topic) == 0 {
		panic("Cannot have empty topic")
	}

	db, err := geoip2.Open("./data/geoip/GeoLite2-Country.mmdb")
	if err != nil {
		panic(err)
	}

	var protocol string
	if secure {
		protocol = "https"
	} else {
		protocol = "http"
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          fmt.Sprintf("metadata.%d.%s", port, protocol),
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(err)
	}

	probeOptions := HttpProbeOptions{Timeout: time.Duration(timeout * 1000000)}
	probeOptions.Secure = secure
	probe := NewHttpProbe(probeOptions)

	collector := workerpool.StartDispatcher(concurrency, concurrency, func(_ string, inputIp string) {
		ip := net.ParseIP(inputIp)

		// ignore err
		country, err := db.Country(ip)
		resp, err := probe.Get(inputIp, port)
		if err != nil {
			log.Printf("[error] %s %s\n", inputIp, err)
			return
		}

		// output to stdin for consumption
		if outputraw {
			fmt.Printf("[success] [%s] [%d] [%s] [%s] %s \n", inputIp, resp.StatusCode, country.Country.IsoCode, protocol, resp.Header)
		} else {

			jsonOutput, _ := json.Marshal(JsonOutput{
				Ip:      inputIp,
				Port:    port,
				Status:  resp.StatusCode,
				Protocol: protocol,
				Country: country.Country.IsoCode,
				Headers: resp.Header,
			})

			fmt.Println(string(jsonOutput))
		}
	})

	for {
		log.Print("Reading from kafka")
		msg, err := c.ReadMessage(-1)
		inputIp := string(msg.Value)
		if err != nil {
			continue
		}

		collector.Work <- workerpool.Work{ID: inputIp, Job: inputIp}
	}
}
