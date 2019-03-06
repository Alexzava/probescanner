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
    // Open log file
    captureFile := "cap.log"
    outFile, err := os.OpenFile(captureFile, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        log.Fatal(err)  
    }
    defer outFile.Close()
    log.SetOutput(outFile)

	// Rad args
	if len(os.Args) < 2 {
		log.Fatal("Invalid arguments")
	}
    log.Println(os.Args)

    // Load vendor database
	LoadVendorDatabase()

    // Start http server
    http.HandleFunc("/scan", HTTPHandler)
    go http.ListenAndServe(":8080", nil)
    fmt.Println("Server online")

	if os.Args[1] == "offline" {
		// Open PCAP file
		pcapFile := os.Args[2]
		handle, err := pcap.OpenOffline(pcapFile)
		if err != nil {
			log.Fatal(err)
		}

		// Scan packets and store information in devicesList
		ScanPackets(handle)

		// Print devices info
		for _, d := range devicesList {
			fmt.Printf("Device (MAC): %s\n\tVendor: %s\n\tSignal: %d\n\tTime (Unix): %d\n\n", d.MAC, d.Vendor, d.RSSI, d.DetectedTime)
		}
	} else if os.Args[1] == "live" {
		interfaces := os.Args[2]
		handle, err := pcap.OpenLive(interfaces, 1600, true, pcap.BlockForever)
		if err != nil {
			log.Fatal(err)
		}
		
		// Scan packets and store information in devicesList
		LiveScan(handle)
		fmt.Println("")

		// Print devices info
		for _, d := range devicesList {
			fmt.Printf("Device (MAC): %s\n\tVendor: %s\n\tSignal: %d\n\tTime (Unix): %d\n\n", d.MAC, d.Vendor, d.RSSI, d.DetectedTime)
		}
	}
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

    	device := DeviceInfo{}
    	device.MAC = dot11.Address2.String()
    	device.BSSID = dot11.Address3.String()
    	device.Vendor = GetVendorInfo(device.MAC)
    	device.RSSI = -100

    	radio := packet.Layer(layers.LayerTypeRadioTap)
    	if radio != nil {
    		dot11r, _ := radio.(*layers.RadioTap)
    		device.RSSI = dot11r.DBMAntennaSignal
    	}
    	device.DetectedTime = time.Now().Unix()

    	_, ok := devicesList[device.MAC]
    	if !ok {
            devicesList[device.MAC] = device
            log.Printf("Device (MAC): %s\n\tVendor: %s\n\tSignal: %d\n\tTime (Unix): %d\n\n", device.MAC, device.Vendor, device.RSSI, device.DetectedTime)
            d++
        }
    	fmt.Printf("\rDevices found: %d", d)
    }
}

// Scan packets and store informations in deviceList
func ScanPackets(handle *pcap.Handle) {
	// Set filter
	var filter string = "subtype probe-req"
    err := handle.SetBPFFilter(filter)
    if err != nil {
        log.Fatal(err)
    }

    // Scan packets
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
    	// Parse 802.11 layer
    	layer := packet.Layer(layers.LayerTypeDot11)
    	dot11, _ := layer.(*layers.Dot11)

    	device := DeviceInfo{}
    	device.MAC = dot11.Address2.String()
    	device.BSSID = dot11.Address3.String()
    	device.Vendor = GetVendorInfo(device.MAC)
    	device.RSSI = -100

    	radio := packet.Layer(layers.LayerTypeRadioTap)
    	if radio != nil {
    		dot11r, _ := radio.(*layers.RadioTap)
    		device.RSSI = dot11r.DBMAntennaSignal
    	}
    	device.DetectedTime = time.Now().Unix()

    	_, ok := devicesList[device.MAC]
    	if !ok {
    		devicesList[device.MAC] = device
    	}
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