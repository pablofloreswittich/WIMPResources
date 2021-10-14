package main

import (
	"log"

	"github.com/unpoller/unifi"
)

type paquete struct {
	CPU                string
	Mem                string
	GeneralTemperature string
	Uptime             string
	FanLevel           string
	MemTotal           string
	MemUsed            string
	MemAverage         string
	Ip                 string
	MAC                int
	Model              int
	Name               int
	Overheating        int
	TxBytes            int
	UpTime             int
	Version            int
	//aca falta todas las mac conectadas al switch
}

func main() {
	c := unifi.Config{
		User: "pablofloreswittich@gmail.com",
		Pass: "pABLO1234pABLO",
		URL:  "https://127.0.0.1:8443/",
		// Log with log.Printf or make your own interface that accepts (msg, fmt)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}
	uni, err := unifi.NewUnifi(&c)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln("Error:", err)
	}

	clients, err := uni.GetClients(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	devices, err := uni.GetDevices(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	/* Clietes conectados con Puerto, MAC propia, IP y MAC del switch. */
	log.Println(len(clients), "Clients connected:")

	for i := 0; i < len(clients); i++ {
		log.Println("Puerto", clients[i].SwPort.Val, "Mac cliente", clients[i].Mac,
			"IP Cliente", clients[i].IP, "Mac SW", clients[i].SwMac)
		/* log.Println("Mac cliente", clients[i].SwPort.Val)
		log.Println("IP Cliente", clients[i].SwPort.Val)
		log.Println("Mac SW", clients[i].SwPort.Val) */
	}

	/* Mac AP, Modelo, Puerto Switch, Mac Switch */
	for i := 0; i < len(devices.UAPs); i++ {
		log.Println("MAC AP", devices.UAPs[i].Mac, "Modelo", devices.UAPs[i].Model,
			"Puerto Sw", devices.UAPs[i].LastUplink.UplinkRemotePort, "Mac SW", devices.UAPs[i].LastUplink.UplinkMac)
		/* log.Println("Modelo", devices.UAPs[i].Model)
		log.Println("Puerto Sw", devices.UAPs[i].LastUplink.UplinkRemotePort)
		log.Println("Mac SW", devices.UAPs[i].LastUplink.UplinkMac) */

	}
	log.Println(len(devices.UAPs))

}
