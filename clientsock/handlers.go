package clientsock

import (
	"encoding/json"
	"errors"

	"github.com/AlessandroRuggiero/go-socketcast"
	fauth "github.com/aleruggiero/cancello-cesana/authentication"
)

func (s *ClientSock) messageHandler(c *socketcast.Client, msg []byte) bool {
	c.Pool.Log.Debug(string(msg))
	var data Message
	err := json.Unmarshal(msg, &data)
	if err != nil {
		c.Pool.Log.Error(err)
		return true
	}
	switch data.Type {
	case 1:
		if s.authClient(c, msg) {
			return true
		}
	case 2:
		c.Pool.Log.Debug("apro")
		if s.apri(c, data) {
			return true
		}
	}
	return false
}

func (s *ClientSock) authClient(c *socketcast.Client, msg []byte) bool {
	var authData AuthRequest
	err := json.Unmarshal(msg, &authData)
	if err != nil {
		c.Pool.Log.Error(err)
		return true
	}
	authToken, err := fauth.CeckToken(authData.Token, s.serverToken)
	if err != nil {
		c.Pool.Log.Infof("autenticazione fallita %s", authData.Token)
		return true
	}
	c.Auth.Authenticated = true
	c.Auth.Token = authData.Token
	// q s t
	c.Metadata.Set(AuthDataS, authToken)
	c.Pool.Log.Infof("autentiucato %s", authToken.Data.UserID)
	if s.Callbacks.SendAperture == nil {
		panic("no and aperture callback")
	}
	s.Callbacks.SendAperture(&CallbackData{
		Client: c,
	}, authToken.Data.UserID)
	return false
}
func (s *ClientSock) apri(c *socketcast.Client, msg Message) bool {
	if !c.Auth.Authenticated {
		c.Pool.Log.Infof("%s non autenticato ha provato a aprire il cancello, chiudo la connessione", c.Conn.RemoteAddr())
		return true
	}
	switch msg.Body {
	case cancello:
		s.apricancello(c)
	case cancelletto:
		s.apricancelletto(c)
	}
	return false
}

func (c *ClientSock) BroadcastAperture(aperture []byte, authId string) {
	c.Pool.ForEach(
		func(c *socketcast.Client) (socketcast.Message, bool, error) {
			if data, found := c.Metadata.Get("authData"); found {
				if token, ok := data.(fauth.FaceToken); ok {
					if token.Data.UserID != authId {
						// fmt.Println("Auth diversi: ",)
						return socketcast.Message{}, false, nil
					}
					return socketcast.Message{
						Type: 5,
						Body: string(aperture),
					}, true, nil

				} else {
					return socketcast.Message{}, false, errors.New("no valid auth")
				}
			} else {
				return socketcast.Message{}, false, errors.New("no auth")
			}
		},
		guard,
	)
}
