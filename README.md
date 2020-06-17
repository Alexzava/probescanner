# probescanner
Probe request scanner to find devices around you. Made with Go!

## Notes
This is an old project, it will not be updated.
It should work. I don't have time to try it, let me know.

## Build

go build -o scanner main.go handler.go


## Usage 

sudo ./scanner < interface >


Interface must be in monitor mode!


sudo airmon-ng start < interface >


## API

**Request**


GET http://127.0.0.1:8683/


**Response**


{
	"00:00:00:00:00:00":{"MAC":"00:00:00:00:00:00","BSSID":"ff:ff:ff:ff:ff:ff","Vendor":"Vendor Name","RSSI":-30,"DetectedTime":1553360100},
	"AA:AA:AA:AA:AA:AA":{"MAC":"AA:AA:AA:AA:AA:AA","BSSID":"ff:ff:ff:ff:ff:ff","Vendor":"Vendor Name","RSSI":-30,"DetectedTime":1553360100},
}
