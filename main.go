package main

import (
	"fmt"
    "log"
    "os"
    "strings"
    "bufio"
    "time"
    "net/http"

    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
)

type DeviceInfo struct {
	MAC string
	BSSID string
	Vendor string
	RSSI int8
	DetectedTime int64
}

var (
	vendorDatabase = make(map[string]string)
	devicesList = make(map[string]DeviceInfo)
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage ./scanner <interface>")
    }

    // Load vendor database
	LoadVendorDatabase()

    // Start http server
    http.HandleFunc("/", HTTPHandler)
    go http.ListenAndServe("127.0.0.1:8683", nil)
    log.Println("Server ready")

    log.Println("Scanning...")
	interfaces := os.Args[1]
    handle, err := pcap.OpenLive(interfaces, 1600, true, pcap.BlockForever)
    if err != nil {
        log.Fatal(err)
    }
    
    // Scan packets and store information in devicesList
    LiveScan(handle)
}

// Scan for probe requests
func LiveScan(handle *pcap.Handle) {
	// Set filter
	err := handle.SetBPFFilter("subtype probe-req")
	if err != nil {
		log.Fatal(err)
	}

	// Scan packets
    devices := 0
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {

    	// Parse 802.11 layer
    	layer := packet.Layer(layers.LayerTypeDot11)

    	if layer == nil {
    		continue
    	}

    	dot11, _ := layer.(*layers.Dot11)

        // Set device informations
    	device := DeviceInfo{}
    	device.MAC = dot11.Address2.String()
    	device.BSSID = dot11.Address3.String()
    	device.Vendor = GetVendorInfo(device.MAC)
        device.DetectedTime = time.Now().Unix()
    	device.RSSI = -100

        // Get RSSI
    	radio := packet.Layer(layers.LayerTypeRadioTap)
    	if radio != nil {
    		dot11r, _ := radio.(*layers.RadioTap)
    		device.RSSI = dot11r.DBMAntennaSignal
    	}

        // Add device to list
    	dd, ok := devicesList[device.MAC]
    	if !ok {
            devicesList[device.MAC] = device
            devices++
        } else {
            if device.RSSI != dd.RSSI {
                delete(devicesList, device.MAC)
                devicesList[device.MAC] = device
            }
        }
    	fmt.Printf("\rDevices found: %d", devices)
    }
}

// Load vendor database
func LoadVendorDatabase() {
	file, err := os.Open("mac.list")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        splitted := strings.Split(line, "\t")
        if len(splitted) > 1 {
            vendorInfo := strings.Replace(line, splitted[0], "", -1)
            vendorInfo = strings.Replace(vendorInfo, "\t", " ", -1)
            vendorDatabase[splitted[0]] = vendorInfo[1:]
        }
    }
}

// Get vendor name by MAC
func GetVendorInfo(MAC string) string {
	splitted := strings.Split(MAC, ":")
    if len(splitted) < 3 {
        return "Unknown"
    }

    vendorMAC := splitted[0] + ":" + splitted[1] + ":" + splitted[2]

    vendorInfo, ok := vendorDatabase[strings.ToUpper(vendorMAC)]
    if ok {
        return vendorInfo
    } else {
        return "Unknown"
    }
}