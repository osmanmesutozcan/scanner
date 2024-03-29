package main

import (
	"io"
	"net"
	"scanner_data_types"
	"strconv"
	"strings"
	"time"
)

func getBannerMessageData(port int) []byte {
	return MessageData[PortMappings[port]]
}

func getBannerProtocolProbedProtocol(port int) string {
	return PortMappings[port]
}

// Read addresses from addrChan and grab banners from these hosts.
// Sends resultStructs to resultChan.  Writes to doneChan when complete.
func GrabBanners(addrChan chan scannertypes.JsonRawIpPort, resultChan chan ProbeResult, doneChan chan int) {
	for addr := range addrChan {
		protocol := getBannerProtocolProbedProtocol(addr.Port)
		deadline := time.Now().Add(time.Duration(*timeoutFlag) * time.Second)
		dialer := net.Dialer{Deadline: deadline}
		conn, err := dialer.Dial("tcp", net.JoinHostPort(addr.Ip, strconv.Itoa(addr.Port)))
		if err != nil {
			resultChan <- ProbeResult{addr.Ip, addr.Port, protocol, nil, err}
			continue
		}

		conn.SetDeadline(deadline)
		s := strings.Replace(string(getBannerMessageData(addr.Port)), "%s", addr.Ip, -1)
		offset := 0
		var buf [1024]byte

		var connectionError error
		for _, line := range strings.Split(s, "##WAIT_ANSWER##\n") {
			if _, err := conn.Write([]byte(line)); err != nil {
				connectionError = err
				break
			}

			n, err := conn.Read(buf[offset:])
			if err != nil && (err != io.EOF || offset == 0) {
				connectionError = err
				break
			}

			offset += n
		}

		if connectionError != nil {
			conn.Close()
			resultChan <- ProbeResult{addr.Ip, addr.Port, protocol, nil, connectionError}
			continue
		}

		conn.Close()
		resultChan <- ProbeResult{addr.Ip, addr.Port, protocol, buf[0:offset], nil}
	}
	doneChan <- 1
}
