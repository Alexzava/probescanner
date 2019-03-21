package main

import (
	"fmt"
    "log"
    "os"
    "io"
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
    // Open log file
    captureFile := "logs"
    outFile, err := os.OpenFile(captureFile, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        log.Fatal(err)  
    }
    defer outFile.Close()

    // Set log outputs
    mw := io.MultiWriter(os.Stdout, outFile)
    log.SetOutput(mw)

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

// Scan packet in live mode
func LiveScan(handle *pcap.Handle) {
	// Set filter
	err := handle.SetBPFFilter("subtype probe-req")
	if err != nil {
		log.Fatal(err)
	}

	// Scan packets
    d := 0
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
    	_, ok := devicesList[device.MAC]
    	if !ok {
            devicesList[device.MAC] = device
            d++
        } else {
            if device.RSSI != devicesList[device.MAC].RSSI {
                delete(devicesList, device.MAC)
                devicesList[device.MAC] = device
            }
        }
    	fmt.Printf("\rDevices found: %d", d)
    }
}

// Load vendor database in memory
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