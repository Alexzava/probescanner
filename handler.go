package main

import (
	"fmt"
	"net/http"
	"io"
)

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	// Send devices info
	var out string
	for _, d := range devicesList {
		out += fmt.Sprintf("Device (MAC): %s\n\tVendor: %s\n\tSignal: %d\n\tTime (Unix): %d\n\n", d.MAC, d.Vendor, d.RSSI, d.DetectedTime)
	}
	io.WriteString(w, out)
}