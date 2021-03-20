package clientsock

import "github.com/AlessandroRuggiero/go-socketcast"

type ClientSock struct {
	Pool      *socketcast.Pool
	Callbacks *Callbacks
}

func (c *ClientSock) Start(callbacks *Callbacks) {
	c.Pool = socketcast.CreatePool(&socketcast.Config{
		OnMessage: c.messageHandler,
	})
	if callbacks == nil {
		callbacks = &Callbacks{}
	}
	if callbacks.Apricancello == nil {
		callbacks.Apricancello = emptycancello
	}
	if callbacks.Apricancelletto == nil {
		callbacks.Apricancelletto = emptycancelletto
	}
	c.Callbacks = callbacks
}

func emptycancello(c *CallbackData) {
	c.Client.Pool.Log.Warn("Nessun callback implementato per apricancello")
}

func emptycancelletto(c *CallbackData) {
	c.Client.Pool.Log.Warn("Nessun callback implementato per apricancelletto")
}

func (c *ClientSock) BroadcastApertoCancello(id string) {
	c.Pool.Broadcast(socketcast.Message{
		Type: 4,
		Body: "cancello",
	}, guard)
}
func (c *ClientSock) BroadcastApertoCancelletto(id string) {
	c.Pool.Broadcast(socketcast.Message{
		Type: 4,
		Body: "cancelletto",
	}, guard)
}

func guard(c *socketcast.Client) bool {
	return c.Auth.Authenticated
}
