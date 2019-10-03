package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"scanner_data_types"
	"syscall"
)

var (
	nConnectFlag = flag.Int("concurrent", 10, "Number of concurrent connections")
	formatFlag   = flag.String("format", "json", "Output format for responses ('ascii', 'hex', json, or 'base64')")
	timeoutFlag  = flag.Int("timeout", 4, "Seconds to wait for each host to respond")
	dataFileFlag = flag.String("data", "", "Directory containing protocol messages to send to responsive hosts ('%s' will be replaced with host IP)")
)

const FileDescLimit = 100000

var MessageData = make(map[string][]byte)
var PortMappings = map[int]string{
	21:   "ftp",
	22:   "ssh",
	80:   "http",
	8080: "http",
}

// Before running main, parse flags and load message data, if applicable
func init() {
	flag.Parse()

	if *dataFileFlag != "" {
		dir, err := ioutil.ReadDir(*dataFileFlag)
		if err != nil {
			panic(err)
		}

		for _, dataFile := range dir {
			dataFileName := dataFile.Name()
			fi, err := os.Open(path.Join(*dataFileFlag, dataFileName))
			if err != nil {
				panic(err)
			}

			buf := make([]byte, 1024)
			n, err := fi.Read(buf)
			MessageData[dataFileName] = buf[0:n]
			if err != nil && err != io.EOF {
				panic(err)
			}
			_ = fi.Close()
		}
	}

	// Increase file descriptor limit
	rlimit := syscall.Rlimit{Max: uint64(FileDescLimit), Cur: uint64(FileDescLimit)}
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[error] cannot set rlimit: %s", err)
	}
}

type ProbeResult struct {
	Addr     string // address of remote host
	Port     int    // connected port of remote host
	Protocol string // probed protocol
	Data     []byte // data returned from the host, if successful
	Err      error  // error, if any
}

func main() {
	addrChan := make(chan scannertypes.JsonRawIpPort, *nConnectFlag) // pass addresses to grabbers
	resultChan := make(chan ProbeResult, *nConnectFlag)              // grabbers send results to output
	doneChan := make(chan int, *nConnectFlag)                        // let grabbers signal completion

	// Start grabbers and output thread
	go Output(resultChan, doneChan)
	for i := 0; i < *nConnectFlag; i++ {
		go GrabBanners(addrChan, resultChan, doneChan)
	}

	// Read addresses from stdin and pass to grabbers
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		addr, err := decodeJson(text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[error] cannot decode payload %s\n", text)
			continue
		}

		addrChan <- addr
	}
	close(addrChan)

	// Wait for completion
	for i := 0; i < *nConnectFlag; i++ {
		<-doneChan
	}
	close(resultChan)
	<-doneChan
}
