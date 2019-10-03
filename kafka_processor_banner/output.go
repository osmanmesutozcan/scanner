package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type JsonOutput struct {
	Addr     string
	Port     int
	Protocol string
	Data     string
	Err      error
}

func probeResultToJsonString(result ProbeResult) string {
	js, err := json.Marshal(JsonOutput{
		result.Addr,
		result.Port,
		result.Protocol,
		string(result.Data),
		result.Err,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] cannot json %s: %s\n", result.Addr, result.Err)
		return ""
	}

	return string(js) + "\n"
}

// Read resultStructs from resultChan, print output, and maintain
// status counters.  Writes to doneChan when complete.
func Output(resultChan chan ProbeResult, doneChan chan int) {
	ok, timeouts, errors := 0, 0, 0

	for result := range resultChan {
		var output string

		switch *formatFlag {
		case "hex":
			output = fmt.Sprintf("%s: %s\n", result.Addr,
				hex.EncodeToString(result.Data))
		case "base64":
			output = fmt.Sprintf("%s: %s\n", result.Addr,
				base64.StdEncoding.EncodeToString(result.Data))
		case "ascii":
			output = fmt.Sprintf("%s: %s\n", result.Addr,
				string(result.Data))
		default:
			output = fmt.Sprintf(probeResultToJsonString(result))
		}

		if result.Err == nil {
			fmt.Printf(output)
			ok++
		} else if nerr, ok := result.Err.(net.Error); ok && nerr.Timeout() {
			fmt.Fprintf(os.Stderr, output)
			timeouts++
		} else {
			fmt.Fprintf(os.Stderr, output)
			errors++
		}
	}

	fmt.Fprintf(os.Stderr, "Complete (OK=%d, timeouts=%d, errors=%d)\n",
		ok, timeouts, errors)

	doneChan <- 1
}
