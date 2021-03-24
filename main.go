package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/AlessandroRuggiero/go-socketcast"
	fauth "github.com/aleruggiero/cancello-cesana/authentication"
	"github.com/aleruggiero/cancello-cesana/clientsock"
	"github.com/aleruggiero/cancello-cesana/database"
	"github.com/aleruggiero/cancello-cesana/espconn"
)

const (
	maxApertureCancello    = 10
	maxApertureCancelletto = 20
)

var db database.CesanaDb
var esp espconn.EspConnection
var clientServer clientsock.ClientSock

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage! %s", os.Getenv("esppassword"))
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests(port string) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/esp", esp.WsEndpoint)
	http.HandleFunc("/ws", clientServer.Pool.Serve)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func comunicaApricencello(c *clientsock.CallbackData) {
	if data, found := c.Client.Metadata.Get(clientsock.AuthDataS); found {
		if token, ok := data.(fauth.FaceToken); ok {
			n := db.GetApertureCancello(token.Data.UserID)
			fmt.Println("Aperture contate", n)
			if n >= maxApertureCancello {
				fmt.Println("Numero aperture massimo raggiunto")
				return
			}
			reqid := db.AddAperturaCancelloAttempt(token.Data.UserID)
			if err := esp.Apricancello(fmt.Sprint(reqid)); err != nil {
				c.Client.Pool.Log.Error(err)
			}
		} else {
			c.Client.Pool.Log.Error("Impossible to cast authdata do facetocken")
		}
	} else {
		c.Client.Pool.Log.Error("Client has no authdata")
	}
}

func comunicaApricencelletto(c *clientsock.CallbackData) {
	fmt.Println("devo aprire il cancelletto")
	if data, found := c.Client.Metadata.Get(clientsock.AuthDataS); found {
		if token, ok := data.(fauth.FaceToken); ok {
			n := db.GetApertureCancelletto(token.Data.UserID)
			if n >= maxApertureCancelletto {
				fmt.Println("Numero aperture massimo raggiunto")
				return
			}
			id := db.AddAperturaCancellettoAttempt(token.Data.UserID)
			if err := esp.Apricancelletto(fmt.Sprint(id)); err != nil {
				c.Client.Pool.Log.Error(err)
			}
		}
	}
}

func aperturaHandler(data string) {
	var authId string
	fmt.Println("data:", data)
	d := strings.Split(data, " ")
	event := d[0]
	id := d[1]
	fmt.Println("hand: ", id)
	switch event {
	case "0":
		authId = db.AddAperturaCancello(id)
		clientServer.BroadcastApertoCancello(id)
	case "1":
		authId = db.AddAperturaCancelletto(id)
		clientServer.BroadcastApertoCancelletto(id)
	default:
		fmt.Println("nnocpt")
		return
	}
	aperture := db.GetAperture(authId)
	ap, err := json.Marshal(&aperture)
	if err != nil {
		fmt.Println("errore parsare aperture")
		return
	}
	clientServer.BroadcastAperture(ap, authId)
}

func sendAperture(c *clientsock.CallbackData, id string) {
	aperture := db.GetAperture(id)
	ap, err := json.Marshal(&aperture)
	if err != nil {
		fmt.Println("errore parsare aperture")
		return
	}
	c.Client.Send(&socketcast.Message{
		Type: 5,
		Body: string(ap),
	})
}

func ceckenv() map[string]string {
	out := make(map[string]string)
	sources := []string{"esppassword", "servertoken", "PORT", "dburi"}
	for _, el := range sources {
		value := os.Getenv(el)
		if len(value) == 0 {
			panic(fmt.Sprintf("manca la variabile env:%s", el))
		}
		out[el] = value
	}
	return out
}

func main() {
	envs := ceckenv()
	err := db.Connect(envs["dburi"])
	if err != nil {
		panic(err)
	}
	defer db.Disconnect()
	clientServer.Start(&clientsock.Callbacks{
		Apricancello:    comunicaApricencello,
		Apricancelletto: comunicaApricencelletto,
		SendAperture:    sendAperture,
	}, envs["servertoken"])
	esp.Password = envs["esppassword"]
	esp.OnAperturaSuccess = aperturaHandler
	handleRequests(envs["PORT"])
}
