package main

import (
	"log"
	"net/http"
	"io"
	"encoding/json"
)

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	//(w).Header().Set("Access-Control-Allow-Origin", "*")
	
	// Create JSON
	out, err := json.Marshal(devicesList)
	if err != nil {
		log.Println(err)
		return
	}

	// Send response to client
	_, err = io.WriteString(w, string(out))
	if err != nil {
		log.Fatal(err)
	}
}