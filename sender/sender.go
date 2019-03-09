/*
    Intercept probe response and fake it.
    Code is not working.

*/

package main

import (
    "fmt"
    "log"
    //"io/ioutil"
    //"os"
    "strings"
    "encoding/hex"

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
    
func main() {
    handle, err := pcap.OpenLive("wlp2s0mon", 1600, true, pcap.BlockForever)
    if err != nil {
        log.Fatal(err)
    }
    
    // Scan packets and store information in devicesList
    PacketFakerLive(handle)
}

func packetkiller() {
    /*buffer = gopacket.NewSerializeBuffer()
    gopacket.SerializeLayers(buffer, options,
        &layers.Ethernet{},
        &layers.IPv4{},
        &layers.TCP{},
        gopacket.Payload(rawBytes),
    )
    outgoingPacket := buffer.Bytes()*/
}

func PacketFakerLive(handle *pcap.Handle) {
    // First Intercept Packet

    // Set filter
    var filter string = "subtype probe-resp"
    err := handle.SetBPFFilter(filter)
    if err != nil {
        log.Fatal(err)
    }

    // Scan packets
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        // Check and Fake Packet

        packetString := fmt.Sprintf("%x", packet.Data())

        //fmt.Printf("%s\n\n", packetString)

        // Change BSSID
        layer := packet.Layer(layers.LayerTypeDot11)
        if layer == nil {
            continue
        }

        dot11, _ := layer.(*layers.Dot11)

        fmt.Printf("Packet from %s to %s\n", dot11.Address2.String(), dot11.Address1.String())

        //dot11.Address2 = []byte{0xf8, 0xdb, 0x7f, 0x7d, 0x4f, 0xbf}
        //dot11.Address3 = []byte{0xf8, 0xdb, 0x7f, 0x7d, 0x4f, 0xbf}
        sa := dot11.Address2.String()
        sa = strings.Replace(sa, ":", "", -1)
        

        if strings.Contains(packetString, sa) {
            fmt.Printf("FoundSA %s\n", sa)
            packetString = strings.Replace(packetString, sa, fmt.Sprintf("%x",[]byte{0xf8, 0xdb, 0x7f, 0x7d, 0x4f, 0xbf}),-1)
        }

        // Change SSID Name
        layer = packet.Layer(layers.LayerTypeDot11InformationElement)
        if layer == nil {
            continue
        }

        dot11Info, _ := layer.(*layers.Dot11InformationElement)
        if dot11Info.ID == layers.Dot11InformationElementIDSSID {
            ssid := fmt.Sprintf("%x", dot11Info.Info)
            newSSID := fmt.Sprintf("%x", []byte("PP WiFi"))
            for len(newSSID) != len(ssid) {
                newSSID += fmt.Sprintf("%x", 0x20)
            }

            if strings.Contains(packetString, ssid) {
                fmt.Printf("FoundSS %s\n", ssid)
                fmt.Printf("NEWSS %s\n", newSSID)
                packetString = strings.Replace(packetString, ssid, newSSID, -1)
            }

            dot11Info.Info = []byte("PP WiFi")
        }

        //ioutil.WriteFile("fake_probe_response.txt", packet.Data(), os.ModeAppend)
        /*layer = packet.Layer(layers.LayerTypeRadioTap)
        if layer == nil {
            continue
        }
        radiotap, _ := layer.(*layers.RadioTap)

        layer = packet.Layer(layers.LayerTypeDot11MgmtProbeResp)
        if layer == nil {
            continue
        }
        dot11MgmtProbeResp, _ := layer.(*layers.Dot11MgmtProbeResp)

        buffer := gopacket.NewSerializeBuffer()
        gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{},
            radiotap,
            dot11,
            dot11MgmtProbeResp,
            dot11Info,
        )
        outgoingPacket := buffer.Bytes()*/

        outgoingPacket, err := hex.DecodeString(packetString)
        if err != nil {
            log.Fatal(err)
        }

        // Send Fake Packet
        err = handle.WritePacketData(outgoingPacket)
        if err != nil {
            log.Fatal(err)
        }

        //fmt.Printf("%x\n", outgoingPacket)

        fmt.Printf("Fake packet from %s to %s\n", dot11.Address2.String(), dot11.Address1.String())
    }
}

func OfflinePacketFaker() {
    // Open PCAP file
    pcapFile := "/home/alex/resp.pcap"
    handle, err := pcap.OpenOffline(pcapFile)
    if err != nil {
        log.Fatal(err)
    }

    // Set filter
    var filter string = "subtype probe-resp"
    err = handle.SetBPFFilter(filter)
    if err != nil {
        log.Fatal(err)
    }

    // Scan packets
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        //fmt.Println(packet)

        // Change BSSID
        layer := packet.Layer(layers.LayerTypeDot11)
        if layer == nil {
            continue
        }

        dot11, _ := layer.(*layers.Dot11)

        fmt.Println(dot11.Address1.String())
        fmt.Println(dot11.Address2.String())
        fmt.Println(dot11.Address3.String())
        fmt.Println("")

        dot11.Address2 = []byte{0xf8, 0xdb, 0x7f, 0x7d, 0x4f, 0xbf}
        dot11.Address3 = []byte{0xf8, 0xdb, 0x7f, 0x7d, 0x4f, 0xbf}

        fmt.Println(dot11.Address1.String())
        fmt.Println(dot11.Address2.String())
        fmt.Println(dot11.Address3.String())
        fmt.Println("")

        // Change SSID Name
        layer = packet.Layer(layers.LayerTypeDot11InformationElement)
        if layer == nil {
            continue
        }

        dot11Info, _ := layer.(*layers.Dot11InformationElement)
        if dot11Info.ID == layers.Dot11InformationElementIDSSID {

            fmt.Println(string(dot11Info.Info))
            dot11Info.Info = []byte("PP WiFi")
            fmt.Println(string(dot11Info.Info))
        }
        break
    }
}