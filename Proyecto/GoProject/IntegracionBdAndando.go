package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	device         string        = "wlo1" //me estoy conectando por la interfaz de wifi de mi computadora
	snapshot_len   int32         = 1024
	promiscuous    bool          = false
	timeout        time.Duration = -1 * time.Second
	pcapFile       string        = "test.pcap"
	handle         *pcap.Handle
	err            error
	snapshotLenuuu uint32 = 1024
	//le puse otra variable por las compatibilidades que presenta con las funciones de abajo
	packetCount int = 0
)

type paquete struct {
	SrcMac    string
	DstMac    string
	ProtoIp   string
	SrcIp     string
	DstIp     string
	ProtoTp   string
	SrcTp     string
	DstTp     string
	ProtoApp  string
	Length    int
	Timestamp time.Time
}

func main() {

	/* Abrimos el dispositivo */
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	/* Seteamos el filtro de captura */
	var filter string = "ether src 60:6D:C7:DF:17:8B"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing MAC SRC 60:6d:c7:df:17:8b packets.")

	/* Iteramos en las capas de los paquetes para sacar informacion. */
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		var p paquete

		/* Acceso a la red */
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethernetLayer != nil {
			ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
			p.SrcMac = ethernetPacket.SrcMAC.String()
			p.DstMac = ethernetPacket.DstMAC.String()
			p.ProtoIp = ethernetPacket.EthernetType.String()
			p.Length = packet.Metadata().CaptureLength
			p.Timestamp = packet.Metadata().Timestamp
		}

		/* Internet */
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			p.SrcIp = ip.SrcIP.String()
			p.DstIp = ip.DstIP.String()
			p.ProtoTp = ip.Protocol.String()
		}

		/* Transporte */
		transportLayer := packet.TransportLayer()
		if transportLayer.LayerType() == layers.LayerTypeTCP {
			tcp, _ := transportLayer.(*layers.TCP)
			p.SrcTp = tcp.SrcPort.String()
			p.DstTp = tcp.DstPort.String()
			p.ProtoApp = tcp.NextLayerType().String()
		}

		if transportLayer.LayerType() == layers.LayerTypeUDP {
			udp, _ := transportLayer.(*layers.UDP)
			p.SrcTp = udp.SrcPort.String()
			p.DstTp = udp.DstPort.String()
			p.ProtoApp = udp.NextLayerType().String()
		}

		/* Configuracion para insertar en la BD */
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://juantuc98:juantuc98@db-wimp.yeslm.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		db := client.Database("wimp")
		col := db.Collection("paquetes")
		if err != nil {
			log.Fatal(err)
		}
		_, err = col.InsertOne(ctx, p)

		packetCount++

	}

}
