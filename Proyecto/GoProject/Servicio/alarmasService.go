package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/unpoller/unifi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type alarma struct {
	Evento    string    `bson:"evento,omitempty,minsize"`
	Mensaje   string    `bson:"mensaje,omitempty,minsize"`
	Timestamp time.Time `bson:"timestamp,omitempty,minsize"`
}

func main() {
	c := unifi.Config{
		User: "pablofloreswittich@gmail.com",
		Pass: "pABLO1234pABLO",
		URL:  "https://127.0.0.1:8443/",
		// Log with log.Printf or make your own interface that accepts (msg, fmt)
		//ErrorLog: log.Printf,
		//DebugLog: log.Printf,
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

		alarmas, err := uni.GetAlarms(sites)
		if err != nil {
			log.Fatalln("Error:", err)
		}

		for i := 0; i < len(alarmas); i++ {
			var a alarma
			a.Evento = alarmas[i].Key
			a.Mensaje = alarmas[i].Msg
			a.Timestamp = alarmas[i].Datetime

			/* Configuracion para insertar en la BD */
			/* client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017")) */
			client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://juantuc98:juantuc98@db-wimp.yeslm.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))

			if err != nil {
				log.Fatal(err)
			}
			/* ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) */
			ctx := context.Background()
			err = client.Connect(ctx)
			db := client.Database("wimp")
			col := db.Collection("alertas")
			if err != nil {
				log.Fatal(err)
			}
			result, err := col.InsertOne(ctx, a)

			fmt.Println(result)
		}
		time.Sleep(60 * time.Second)
	}
}
