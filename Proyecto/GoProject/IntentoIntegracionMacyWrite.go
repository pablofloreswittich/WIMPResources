package main

import (
	"fmt"
	"log"
	"time"
	"os"


	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

var (
	device       string = "wlo1"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	timeout      time.Duration = -1 * time.Second
	pcapFile string = "test.pcap"
	handle   *pcap.Handle
	err      error
	snapshotLen    int32  = 1024
	snapshotLenuuu uint32 = 1024
	//le puse otra variable por las compatibilidades que presenta con las funciones de abajo
	packetCount int = 0
)



func main() {

	// Open output pcap file and write header
	f, _ := os.Create("test.pcap")
	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(snapshotLenuuu, layers.LinkTypeEthernet)
	defer f.Close()

	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	var filter string = "ether src 60:6D:C7:DF:17:8B"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing MAC SRC 60:6d:c7:df:17:8b packets.")

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Do something with a packet here.
		fmt.Println(packet)
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())

		
		packetCount++

		// Only capture 100 and then stop
		if packetCount > 100 {
			break
		}
	}




}
