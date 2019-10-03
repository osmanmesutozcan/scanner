package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"os"
)

// Read resultStructs from resultChan, print output, and maintain
// status counters.  Writes to doneChan when complete.
func output(resultChan chan resultStruct, doneChan chan int) {
	ok, timeouts, errors := 0, 0, 0
	for result := range resultChan {
		if result.err == nil {
			switch *formatFlag {
			case "hex":
				fmt.Printf("%s: %s\n", result.addr,
					hex.EncodeToString(result.data))
			case "base64":
				fmt.Printf("%s: %s\n", result.addr,
					base64.StdEncoding.EncodeToString(result.data))
			default:
				fmt.Printf("%s: %s\n", result.addr,
					string(result.data))
			}
			ok++
		} else if nerr, ok := result.err.(net.Error); ok && nerr.Timeout() {
			fmt.Fprintf(os.Stderr, "%s: Timeout\n", result.addr)
			timeouts++
		} else {
			fmt.Fprintf(os.Stderr, "%s: Error %s\n", result.addr, result.err)
			errors++
		}
	}

	fmt.Fprintf(os.Stderr, "Complete (OK=%d, timeouts=%d, errors=%d)\n",
		ok, timeouts, errors)

	doneChan <- 1
}
