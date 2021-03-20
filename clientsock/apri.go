package clientsock

import "github.com/AlessandroRuggiero/go-socketcast"

func (s *ClientSock) apricancello(c *socketcast.Client) {
	s.Callbacks.Apricancello(&CallbackData{
		Client: c,
	})
}
func (s *ClientSock) apricancelletto(c *socketcast.Client) {
	s.Callbacks.Apricancelletto(&CallbackData{
		Client: c,
	})
}
