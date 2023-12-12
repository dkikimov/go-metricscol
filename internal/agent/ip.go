package agent

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
)

var currentIP net.IP
var err error
var triedToGet atomic.Bool

func getOutboundIP() (net.IP, error) {
	if triedToGet.Load() == false {
		defer triedToGet.Store(true)

		response, err := http.Get("https://eth0.me")
		if err != nil {
			err = errors.New("couldn't reach eth0.me website")
			return nil, err
		}

		responseBytes, err := io.ReadAll(response.Body)
		if err != nil {
			err = errors.New("couldn't read response body")
			return nil, err
		}

		stringIP := string(responseBytes)
		stringIP = strings.Trim(stringIP, "\n")
		parsedIP := net.ParseIP(stringIP)
		if parsedIP == nil {
			err = errors.New("couldn't parse ip address from string")
			return nil, err
		}

		currentIP = parsedIP
	}

	if err != nil {
		return nil, err
	}
	return currentIP, nil
}
