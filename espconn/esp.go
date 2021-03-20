package espconn

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type EspConnection struct {
	Conn              *websocket.Conn
	OnAperturaSuccess func(string)
	Password          string
}

func (esp *EspConnection) reader(conn *websocket.Conn) {
	for {
		// read in a message
		_, p, err := conn.ReadMessage() // _ sarebbe messagetype
		if err != nil {
			log.Println(err)
			log.Println("Errore in reader")
			if websocket.IsCloseError(err) {
				log.Println("Andato viaaaa")
				esp.Conn = nil
			}
			return
		}
		// print out that message for clarity
		msg := string(p)
		log.Println(msg)
		if msg == esp.Password {
			fmt.Println("esp verificato collegato")
			esp.Conn = conn
			continue
		}
		if esp.OnAperturaSuccess == nil {
			panic("no apertura handler")
		}
		esp.OnAperturaSuccess(msg)

	}

}

func (esp *EspConnection) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	esp.reader(ws)
}
func (esp *EspConnection) WriteMessage(messageType int, data []byte) error {
	if !esp.Connected() {
		return errors.New("esp non connesso")
	}
	return esp.Conn.WriteMessage(messageType, data)
}
func (esp *EspConnection) Connected() bool {
	return esp.Conn != nil
}
