package main

import (
	"encoding/json"
	"scanner_data_types"
)

func decodeJson(addressInput string) (scannertypes.JsonRawIpPort, error) {
	var result scannertypes.JsonRawIpPort
	err := json.Unmarshal([]byte(addressInput), &result)
	return result, err
}

