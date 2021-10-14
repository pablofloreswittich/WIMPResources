package main

import (
	"fmt"
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

	alarmas, err := uni.GetAlarms(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	/*
		anos, err := uni.GetAnomalies(sites)
		if err != nil {
			log.Fatalln("Error:", err)
		} */

	/* 	log.Println(devices.USWs[0].SystemStats.CPU, "    CPU")
	   	log.Println(devices.USWs[0].SystemStats.Mem, "    Mem")
	   	log.Println(devices.USWs[0].GeneralTemperature, "    GeneralTemperature")
	   	log.Println(devices.USWs[0].SystemStats.Uptime, "    Uptime")
	   	log.Println(devices.USWs[0].FanLevel, "    FanLevel")
	   	log.Println(devices.USWs[0].SysStats.MemTotal.Val, "    MemTotal")
	   	log.Println(devices.USWs[0].SysStats.MemUsed.Val, "    MemUsed")
	   	log.Println(int((float32(devices.USWs[0].SysStats.MemUsed.Val)/float32(devices.USWs[0].SysStats.MemTotal.Val))*100), "    MemAverage")
	   	log.Println(devices.USWs[0].ID, "    Id")
	   	log.Println(devices.USWs[0].IP, "    Ip")
	   	log.Println(devices.USWs[0].Mac, "    MAC")
	   	log.Println(devices.USWs[0].Model, "    Model")
	   	log.Println(devices.USWs[0].Name, "    Name")
	   	log.Println(devices.USWs[0].Overheating, "    Overheating")
	   	log.Println(devices.USWs[0].TxBytes, "    TxBytes")
	   	log.Println(devices.USWs[0].Uptime, "    UpTime")
	   	log.Println(devices.USWs[0].Version, "    Version") */

	log.Println(len(devices.USGs), "Unifi Gateways Found")

	log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
	for i, uap := range devices.UAPs {
		log.Println(i+1, uap.Name, uap.IP)
	}

	log.Println(alarmas[0])
	log.Println(len(alarmas), "Alarmas. ")
	for i := 0; i < len(alarmas); i++ {
		fmt.Println("")
		log.Println("Evento ", i+1, ": ", alarmas[i].Key)
		log.Println("Mensaje ", i+1, ": ", alarmas[i].Msg)
		log.Println("Fecha ", i+1, ": ", alarmas[i].Datetime)
		fmt.Println("")

	}

	for i := 0; i < len(clients); i++ {
		/* log.Println(clients[i].Mac) */
		log.Println(clients[i].SwPort)
		log.Println(clients[i].Mac)
	}

	for i := 0; i < len(clients); i++ {
		/* log.Println(clients[i].Mac) */
		log.Println(clients[i].SwPort)
		log.Println(clients[i].Mac)
	}

	log.Println("-----------------------------------------------------------------------")

	log.Println(len(clients), "Clients connected:")
	for i, client := range clients {
		log.Println(i+1, client.SwPort.Val, client.Mac, client.IP, client.SwMac)
	}

	/* Mac AP, Modelo, Puerto Switch, Mac Switch */
	log.Println(len(devices.UAPs), "AP Connected:")
	for i, ap := range devices.UAPs {
		log.Println(i+1, ap.Mac, ap.Model, ap.LastUplink.UplinkRemotePort, ap.LastUplink.UplinkMac)
	}

	log.Println("-----------------------------------------------------------------------")

	log.Println(devices.USWs[0].SystemStats.CPU.Val, "    CPU")
	log.Println(devices.USWs[0].GeneralTemperature.Val, "    GeneralTemperature")
	log.Println(devices.USWs[0].SystemStats.Uptime.Val, "    Uptime")
	log.Println(devices.USWs[0].FanLevel.Val, "    FanLevel")
	log.Println(int((float32(devices.USWs[0].SysStats.MemUsed.Val)/float32(devices.USWs[0].SysStats.MemTotal.Val))*100), "    MemAverage")
	log.Println(devices.USWs[0].IP, "    Ip")
	log.Println(devices.USWs[0].Mac, "    MAC")
	log.Println(devices.USWs[0].Model, "    Model")
	log.Println(devices.USWs[0].Name, "    Name")
	log.Println(devices.USWs[0].DownlinkTable, "    Version")

}
