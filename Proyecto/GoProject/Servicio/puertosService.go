package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/unpoller/unifi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Qty struct {
	NumClients int
	NumAp      int
}

type Ports struct {
	Ports []InfoPort
}

type InfoPort struct {
	Num        int       `bson:"num,omitempty,minsize"`
	Mac        string    `bson:"mac,omitempty,minsize"`
	Ip         string    `bson:"ip,omitempty,minsize"`
	Name       string    `bson:"name,omitempty,minsize"`
	Model      string    `bson:"model,omitempty,minsize"`
	Uptime     int       `bson:"uptime,omitempty,minsize"`
	CPU        int       `bson:"cpu,omitempty,minsize"`
	Mem        int       `bson:"mem,omitempty,minsize"`
	ClientesAp []Puertos `bson:"clientesap,omitempty,minsize"`
}

type Puertos struct {
	Mac string `bson:"mac,omitempty,minsize"`
	Ip  string `bson:"ip,omitempty,minsize"`
}

func main() {
	c := unifi.Config{
		User: "pablofloreswittich@gmail.com",
		Pass: "pABLO1234pABLO",
		URL:  "https://127.0.0.1:8443/",
		// Log with log.Printf or make your own interface that accepts (msg, fmt)
		/* ErrorLog: log.Printf, */
		/* DebugLog: log.Printf, */
	}
	uni, err := unifi.NewUnifi(&c)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	for {
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
		var Indice int
		var arrSw Ports
		var filtro string
		var qty Qty
		qty.NumClients = len(clients)
		qty.NumAp = len(devices.UAPs)
		//switches := devices.USWs

		//	for y, switch := range switches {
		// macswitchactual = switch.mac

		for i := 0; i < qty.NumAp; i++ {
			//if macswitchactual == devices.UAPs[i].LastUplink.UplinkMac (mac del sw proximo)
			var info InfoPort
			info.Mac = devices.UAPs[i].Mac
			info.Ip = devices.UAPs[i].IP
			info.Name = devices.UAPs[i].Name
			info.Model = devices.UAPs[i].Model
			info.Uptime = int(devices.UAPs[i].Uptime.Val)
			info.CPU = int(devices.UAPs[i].SystemStats.CPU.Val)
			info.Mem = int((float32(devices.UAPs[i].SysStats.MemUsed.Val) / float32(devices.UAPs[i].SysStats.MemTotal.Val)) * 100)
			info.Num = int(devices.UAPs[i].LastUplink.UplinkRemotePort)
			/* 		fmt.Println("Mac arriba de AP", devices.UAPs[i].LastUplink.UplinkMac) */
			arrSw.Ports = append(arrSw.Ports, info)
		}

		for i := 0; i < qty.NumClients; i++ {
			var info InfoPort
			var info2 Puertos
			filtro = clients[i].SwMac //mac sw proximo
			//if macsw == filtro

			// fmt.Println(clients[i].ApMac) MAC DEL AP ARRIBA MIO, SI ES QUE HAY.
			// fmt.Println(clients[i].SwMac) MAC DEL PRIMER SW DEL ARBOL.
			info.Mac = clients[i].Mac
			info.Ip = clients[i].IP
			info.Num = int(clients[i].SwPort.Val)
			info2.Mac = clients[i].Mac
			info2.Ip = clients[i].IP

			if clients[i].ApMac != "" {
				for u, elemento := range arrSw.Ports {
					if clients[i].ApMac == elemento.Mac {
						Indice = u
						break
					}
				}

				arrSw.Ports[Indice].ClientesAp = append(arrSw.Ports[Indice].ClientesAp, info2)
			} else {
				// si no esta vinculado a un AP -> esta vinculado a un sw
				//if macswitchactual == clients[i].SwMac (mac del sw proximo)

				arrSw.Ports = append(arrSw.Ports, info)
				//else
				//continue
			}
			//

		}
		//} corta el lazo switches
		fmt.Println(arrSw.Ports)

		/* Configuracion para insertar en la BD */
		//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://juantuc98:juantuc98@db-wimp.yeslm.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		err = client.Connect(ctx)
		db := client.Database("wimp")
		col := db.Collection("switches")
		opts := options.Update().SetUpsert(true)
		filter := bson.D{{"mac", filtro}}
		update := bson.D{
			{"$set",
				bson.D{
					{"clients", qty.NumClients},
					{"aps", qty.NumAp},
					{"ports", arrSw.Ports},
				},
			},
		}
		if err != nil {
			log.Fatal(err)
		}
		result, err := col.UpdateOne(ctx, filter, update, opts)
		/* fmt.Println(err) */
		fmt.Println(result)

		time.Sleep(60 * time.Second)
	}
}
