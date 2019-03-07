package main

import (
	//"fmt"
	"log"
	"net/http"
	//"io"
	"encoding/json"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {return true}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	
	if string(message) == "get_devices" {
		// Send devices info
		out, err := json.Marshal(devicesList)
		if err != nil {
			log.Println(err)
			return
		}
		err = c.WriteMessage(mt, out)
		if err != nil {
			log.Println(err)
			return
		}
	}
}